
const chai = require('chai');
const  { expect } = chai;
const helper = require('./helper')

describe('The "bridg-admin" role', () => {

  const sub = 'bridg-admin';
  let action;
  let resourceName;
  let response;

  const updateWardenResponse = (act, rn, done) => {
    action = act;
    resourceName = rn;
    helper.makeWardenReq(sub, action, resourceName, null, (err, res) => {
      response = res;
      done();
    });
  };

  describe('HTTP "GET"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/brands"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/brands"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/sites"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/sites"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/roles"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:roles', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/roles"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:roles', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/users"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:users', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/users"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:users', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/audiences"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/aud_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:audiences:aud_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/aud_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:audiences:aud_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/analytics/campaigns"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:analytics:campaigns', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/analytics/campaigns"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:analytics:campaigns', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences/aud_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences/aud_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshots"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots/snp_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots:snp_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshots/snp_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots:snp_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports/exp_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports:exp_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports/exp_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports:exp_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/client-configuration"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:client-configuration', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/client-configuration"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_1:client-configuration', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/insights"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:insights', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/insights"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:insights', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/reveal-jobs"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/reveal-jobs"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/reveal-jobs/1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/reveal-jobs/2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/reveal-jobs/latest"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:latest', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/reveal-jobs/latest"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:latest', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/reveal-jobs/latest/artifact"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:latest:artifact', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/reveal-jobs/latest/artifact"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:latest:artifact', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/sites"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/sites"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/sites/ste_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/sites/ste_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/snapshot-fb-export-facebook-account"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/snapshot-fb-export-facebook-account"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/campaigns"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:campaigns', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/clt_1"', (done) => {
        updateWardenResponse('read', 'rn:hydra:clients:clt_1', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

    describe('context "/integrations"', () => {
      it('is allowed access to "/int_1/sync-agent-instances"', (done) => {
        updateWardenResponse('read', 'rn:bridg:integrations:int_1:sync-agent-instances', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:hydra:policies', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/pol_1"', (done) => {
        updateWardenResponse('read', 'rn:hydra:policies:pol_1', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

    describe('context "/reveal-jobs"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:reveal-jobs', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rvl_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:reveal-jobs:rvl_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rvl_1/artifact"', (done) => {
        updateWardenResponse('read', 'rn:bridg:reveal-jobs:rvl_1:artifact', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/roles"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:roles', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:roles:rle_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_1/users"', (done) => {
        updateWardenResponse('read', 'rn:bridg:roles:rle_1:users', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/search"', () => {
      it('is allowed access to "/crm-txn/customer-profile/_search"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:crm-txn:customer-profile:_search', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/crm-txn/customer-profile/_count"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:crm-txn:customer-profile:_count', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/customer-profile/_search"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:act_1:customer-profile:_search', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/customer-profile/_search"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:act_2:customer-profile:_search', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/customer-profile/_count"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:act_1:customer-profile:_count', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/customer-profile/_count"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:act_2:customer-profile:_count', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is NOT allowed access to "/crm-txn/_search"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:crm-txn:_search', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/crm-txn/_count"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:crm-txn:_count', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/act_1/_search"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:act_1:_search', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/act_2/_search"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:act_1:_search', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/act_1/_count"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:act_1:_count', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/act_2/_count"', (done) => {
        updateWardenResponse('read', 'rn:bridg:search:act_2:_count', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/ste_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/ste_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/snapshot-fb-export-facebook-accounts"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:snapshot-fb-export-facebook-accounts', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/snapshot-fb-exports"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:snapshot-fb-exports', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/accounts"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_1:accounts', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/accounts"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_2:accounts', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/authorizations"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_1:authorizations', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/authorizations"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_2:authorizations', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/brands"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_1:brands', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/brands"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_2:brands', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/roles"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_1:roles', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/roles"', (done) => {
        updateWardenResponse('read', 'rn:bridg:users:usr_2:roles', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

  });

  describe('HTTP "POST"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/actions/add-user"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_1:actions:add-user', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/actions/add-user"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_2:actions:add-user', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/users/usr_1"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_1:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/users/usr_1"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_2:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/analytics"', () => {
      it('is allowed access to "/campaigns/revenue"', (done) => {
        updateWardenResponse('create', 'rn:bridg:analytics:campaigns:revenue', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/authenticate"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:bridg:authenticate', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is allowed access to "/brd_1/analytics/campaigns/bychannel/facebook"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:analytics:campaigns:bychannel:facebook', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/analytics/campaigns/bychannel/facebook"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:analytics:campaigns:bychannel:facebook', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshots', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/snapshot-fb-export-facebook-account"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:snapshot-fb-export-facebook-account', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/snapshot-fb-export-facebook-account"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:snapshot-fb-export-facebook-account', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/sites"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/sites', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:hydra:clients', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

    describe('context "/integrations"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:bridg:integrations', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/int_1/actions/regenerate-access-key"', (done) => {
        updateWardenResponse('create', 'rn:bridg:integrations:int_1:actions:regenerate-access-key', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/integration-types"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:bridg:integration-types', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/metrics"', () => {
      it('is allowed access to "/write"', (done) => {
        updateWardenResponse('create', 'rn:bridg:metrics:write', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/query"', (done) => {
        updateWardenResponse('create', 'rn:bridg:metrics:query', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:hydra:policies', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

    describe('context "/reveal-jobs"', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:bridg:reveal-jobs', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/roles" (rle_1: bridg-admin, rle_2: account-admin, usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:bridg:roles', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_1/users/usr_1"', (done) => {
        updateWardenResponse('create', 'rn:bridg:roles:rle_1:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_1/users/usr_2"', (done) => {
        updateWardenResponse('create', 'rn:bridg:roles:rle_1:users:usr_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_2/users/usr_1"', (done) => {
        updateWardenResponse('create', 'rn:bridg:roles:rle_2:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_2/users/usr_2"', (done) => {
        updateWardenResponse('create', 'rn:bridg:roles:rle_2:users:usr_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/scheduler"', () => {
      it('is allowed access to "/run"', (done) => {
        updateWardenResponse('create', 'rn:bridg:scheduler:run', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/schedule"', (done) => {
        updateWardenResponse('create', 'rn:bridg:scheduler:schedule', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/search"', () => {
      it('is allowed access to "/crm-txn/customer-profile/_search"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:crm-txn:customer-profile:_search', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/crm-txn/customer-profile/_count"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:crm-txn:customer-profile:_count', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/customer-profile/_search"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:act_1:customer-profile:_search', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/customer-profile/_search"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:act_2:customer-profile:_search', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/customer-profile/_count"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:act_1:customer-profile:_count', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/customer-profile/_count"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:act_2:customer-profile:_count', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is NOT allowed access to "/crm-txn/_search"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:crm-txn:_search', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/crm-txn/_count"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:crm-txn:_count', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/act_1/_search"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:act_1:_search', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/act_2/_search"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:act_1:_search', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/act_1/_count"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:act_1:_count', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/act_2/_count"', (done) => {
        updateWardenResponse('create', 'rn:bridg:search:act_2:_count', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:bridg:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/ste_1/location-links"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1:location-links', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/ste_2/location-links"', (done) => {
        updateWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_1:location-links', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/"', (done) => {
        updateWardenResponse('create', 'rn:bridg:users', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/actions/confirm-account-activation"', (done) => {
        updateWardenResponse('create', 'rn:bridg:users:usr_1:actions:confirm-account-activation', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/actions/confirm-account-activation', (done) => {
        updateWardenResponse('create', 'rn:bridg:users:usr_2:actions:confirm-account-activation', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/actions/confirm-password-reset"', (done) => {
        updateWardenResponse('create', 'rn:bridg:users:usr_1:actions:confirm-password-reset', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/actions/confirm-password-reset"', (done) => {
        updateWardenResponse('create', 'rn:bridg:users:usr_2:actions:confirm-password-reset', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/actions/request-password-reset"', (done) => {
        updateWardenResponse('create', 'rn:bridg:users:usr_1:actions:request-password-reset', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/actions/request-password-reset"', (done) => {
        updateWardenResponse('create', 'rn:bridg:users:usr_2:actions:request-password-reset', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/warden"', () => {
      it('is NOT allowed access to "/allowed"', (done) => {
        updateWardenResponse('decide', 'rn:hydra:warden:allowed', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });

      it('is NOT allowed access to "/token/allowed"', (done) => {
        updateWardenResponse('decide', 'rn:hydra:warden:token:allowed', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

  });

  describe('HTTP "PUT"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/act_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/actions/edit-user"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_1:actions:edit-user', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/actions/edit-user"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_2:actions:edit-user', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is allowed access to "/brd_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences/aud_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences/aud_2"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports/exp_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports:exp_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports/exp_2"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports:exp_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/integration-types"', () => {
      it('is allowed access to "/it_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:integration-types:it_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/roles"', () => {
      it('is allowed access to "/rle_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:roles:rle_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is allowed access to "/ste_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/ste_2"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/ste_1/location-links/ll_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1:location-links:ll_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/ste_2/location-links/ll_2"', (done) => {
        updateWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_1:location-links:ll_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/usr_1"', (done) => {
        updateWardenResponse('update', 'rn:bridg:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2"', (done) => {
        updateWardenResponse('update', 'rn:bridg:users:usr_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/change-email"', (done) => {
        updateWardenResponse('update', 'rn:bridg:users:usr_1:change-email', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/change-email"', (done) => {
        updateWardenResponse('update', 'rn:bridg:users:usr_2:change-email', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_1/change-password"', (done) => {
        updateWardenResponse('update', 'rn:bridg:users:usr_1:change-password', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/usr_2/change-password"', (done) => {
        updateWardenResponse('update', 'rn:bridg:users:usr_2:change-password', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

  });

  describe('HTTP "DELETE"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/act_1/actions/remove-user"', (done) => {
        updateWardenResponse('delete', 'rn:bridg:accounts:act_1:actions:remove-user', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/actions/remove-user"', (done) => {
        updateWardenResponse('delete', 'rn:bridg:accounts:act_2:actions:remove-user', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_1/users/usr_1"', (done) => {
        updateWardenResponse('delete', 'rn:bridg:accounts:act_1:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/act_2/users/usr_2"', (done) => {
        updateWardenResponse('delete', 'rn:bridg:accounts:act_2:users:usr_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/clt_1"', (done) => {
        updateWardenResponse('delete', 'rn:hydra:clients:clt_1', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/pol_1"', (done) => {
        updateWardenResponse('delete', 'rn:hydra:policies:pol_1', () => {
          expect(response.body.allowed).to.equal(false);
          done();
        });
      });
    });

    describe('context "/roles" (rle_1: bridg-admin, rle_2: account-admin, usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/rle_1/users/usr_1"', (done) => {
        updateWardenResponse('delete', 'rn:bridg:roles:rle_1:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_1/users/usr_2"', (done) => {
        updateWardenResponse('delete', 'rn:bridg:roles:rle_1:users:usr_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_2/users/usr_1"', (done) => {
        updateWardenResponse('delete', 'rn:bridg:roles:rle_2:users:usr_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/rle_2/users/usr_2"', (done) => {
        updateWardenResponse('delete', 'rn:bridg:roles:rle_2:users:usr_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

  });

});
