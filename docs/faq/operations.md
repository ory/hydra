# Operational Considerations

This section is intended to give operations folk some useful guidelines on
how to best manage a hydra deployment.

Note: some of this section is specific to rethinkdb, which is the only persistence
engine supported at the time of writing.

## Managing Client/Policy Definitions

It is useful for JSON files for client and policy definitions to be persisted outside
of the hydra database itself. There are several reasons for this:

- client secrets are stored in the database as a bcrypt hash, so you will not be
  able to read back the secret.  If, at some later point, you need to update the
  client definition then you need to delete the client definition from the database
  and recreate it. So as to ensure that existing web/mobile/other apps are able to
  continue operation, you will need to recreate the same client ID/secret - for which
  you should use the `hydra clients import` command.

A good storage platform for these files is Hashicorp [Vault](https://www.vaultproject.io). The full setup of a
Vault deployment is beyond the scope of this document, please refer to the Vault website
for this.

With LDAP-based authentication configured, reading a client definition from
the Vault "secrets" backend is as simple as:

````
$ vault auth -method=ldap username=john.doe
$ vault read -field=myclient secrets/hydra/clients | hydra clients import /dev/stdin
````

The first command above will prompt you for your LDAP password and then create
a file `~/.vault-token`. Note that you may wish to consider using token authentication
instead of LDAP, but the above should be simple for ops folk to perform.

If you need to update a client later, you can delete the client from hydra using:

````
$ hydra clients delete "<clientid>"
````
and then re-import as above.

The same procedure can be followed for importing/updating policy definitions.

## Recovering root client access

If you somehow manage to lose admin access to your Hydra system, you can regain this
by making use of Hydra's temporary root client creation - which is triggered when
hydra is unable to find any client definitions upon startup. Due to the ID given to
policy used for temporary root clients, you may need to also delete configured
policies. With a rethinkdb connection `r`, you can perform these operations as follows:

````
r.db('hydra').table('hydra_clients').delete()
r.db('hydra').table('hydra_policies').delete()
````

then: 

- restart Hydra
- re-import your client/policy definitions, as described above
- delete your new temporary root client
- ensure that any Hydra clients which have read keys from hydra are refreshed, possibly 
  involving a simple restart to effect a timely update
- play the maracas, FTW
