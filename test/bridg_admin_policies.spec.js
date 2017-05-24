
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

  const testWardenResponse = (act, rn, exp, done) => {
    updateWardenResponse(act, rn, () => {
      expect(response.body.allowed).to.equal(exp);
      done();
    });
  };

  describe('HTTP "GET"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts', true, done);
      });

      it('is allowed access to "/act_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1', true, done);
      });

      it('is allowed access to "/act_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2', true, done);
      });

      it('is allowed access to "/act_1/brands"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands', true, done);
      });

      it('is allowed access to "/act_2/brands"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands', true, done);
      });

      it('is allowed access to "/act_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:search:customer-profile:_search', true, done);
      });

      it('is allowed access to "/act_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:search:customer-profile:_count', true, done);
      });

      it('is allowed access to "/act_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:search:customer-profile:_search', true, done);
      });

      it('is allowed access to "/act_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:search:customer-profile:_count', true, done);
      });

      it('is allowed access to "/act_1/sites"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:sites', true, done);
      });

      it('is allowed access to "/act_2/sites"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:sites', true, done);
      });

      it('is allowed access to "/act_1/roles"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:roles', true, done);
      });

      it('is allowed access to "/act_2/roles"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:roles', true, done);
      });

      it('is allowed access to "/act_1/users"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:users', true, done);
      });

      it('is allowed access to "/act_2/users"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:users', true, done);
      });
    });

    describe('context "/audiences"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:audiences', true, done);
      });

      it('is allowed access to "/aud_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:audiences:aud_1', true, done);
      });

      it('is allowed access to "/aud_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:audiences:aud_2', true, done);
      });
    });

    describe('context "/audience-export-groups"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:audience-export-groups', true, done);
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:brands', true, done);
      });

      it('is allowed access to "/brd_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1', true, done);
      });

      it('is allowed access to "/brd_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2', true, done);
      });

      it('is allowed access to "/brd_1/analytics/campaigns"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:analytics:campaigns', true, done);
      });

      it('is allowed access to "/brd_2/analytics/campaigns"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:analytics:campaigns', true, done);
      });

      it('is allowed access to "/brd_1/audiences"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', true, done);
      });

      it('is allowed access to "/brd_2/audiences"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', true, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1', true, done);
      });

      it('is allowed access to "/brd_2/audiences/aud_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2', true, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots', true, done);
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshots"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots', true, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots/snp_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots:snp_1', true, done);
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshots/snp_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots:snp_2', true, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports', true, done);
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports', true, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports/exp_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports:exp_1', true, done);
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports/exp_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports:exp_2', true, done);
      });

      it('is allowed access to "/brd_1/audience-export-groups"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audience-export-groups', true, done);
      });

      it('is allowed access to "/brd_2/audience-export-groups"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audience-export-groups', true, done);
      });

      it('is allowed access to "/brd_1/audience-export-groups/aeg_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audience-export-groups:aeg_1', true, done);
      });

      it('is allowed access to "/brd_2/audience-export-groups/aeg_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audience-export-groups:aeg_2', true, done);
      });

      it('is allowed access to "/brd_1/client-configuration"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:client-configuration', true, done);
      });

      it('is allowed access to "/brd_2/client-configuration"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:client-configuration', true, done);
      });

      it('is allowed access to "/brd_1/insights"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:insights', true, done);
      });

      it('is allowed access to "/brd_2/insights"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:insights', true, done);
      });

      it('is allowed access to "/brd_1/reveal-jobs"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs', true, done);
      });

      it('is allowed access to "/brd_2/reveal-jobs"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs', true, done);
      });

      it('is allowed access to "/brd_1/reveal-jobs/1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:1', true, done);
      });

      it('is allowed access to "/brd_2/reveal-jobs/2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:2', true, done);
      });

      it('is allowed access to "/brd_1/reveal-jobs/latest"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:latest', true, done);
      });

      it('is allowed access to "/brd_2/reveal-jobs/latest"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:latest', true, done);
      });

      it('is allowed access to "/brd_1/reveal-jobs/latest/artifact"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:latest:artifact', true, done);
      });

      it('is allowed access to "/brd_2/reveal-jobs/latest/artifact"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:latest:artifact', true, done);
      });

      it('is allowed access to "/brd_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_search', true, done);
      });

      it('is allowed access to "/brd_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_count', true, done);
      });

      it('is allowed access to "/brd_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_search', true, done);
      });

      it('is allowed access to "/brd_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_count', true, done);
      });

      it('is allowed access to "/brd_1/sites"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites', true, done);
      });

      it('is allowed access to "/brd_2/sites"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites', true, done);
      });

      it('is allowed access to "/brd_1/sites/ste_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', true, done);
      });

      it('is allowed access to "/brd_2/sites/ste_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', true, done);
      });

      it('is allowed access to "/brd_1/snapshot-fb-export-facebook-account"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', true, done);
      });

      it('is allowed access to "/brd_2/snapshot-fb-export-facebook-account"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', true, done);
      });
    });

    describe('context "/campaigns"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:campaigns', true, done);
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/clt_1"', (done) => {
        testWardenResponse('read', 'rn:hydra:clients:clt_1', false, done);
      });
    });

    describe('context "/integrations"', () => {
      it('is allowed access to "/int_1/sync-agent-instances"', (done) => {
        testWardenResponse('read', 'rn:bridg:integrations:int_1:sync-agent-instances', true, done);
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:hydra:policies', false, done);
      });

      it('is NOT allowed access to "/pol_1"', (done) => {
        testWardenResponse('read', 'rn:hydra:policies:pol_1', false, done);
      });
    });

    describe('context "/reveal-jobs"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:reveal-jobs', true, done);
      });

      it('is allowed access to "/rvl_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:reveal-jobs:rvl_1', true, done);
      });

      it('is allowed access to "/rvl_1/artifact"', (done) => {
        testWardenResponse('read', 'rn:bridg:reveal-jobs:rvl_1:artifact', true, done);
      });
    });

    describe('context "/roles"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:roles', true, done);
      });

      it('is allowed access to "/rle_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:roles:rle_1', true, done);
      });

      it('is allowed access to "/rle_1/users"', (done) => {
        testWardenResponse('read', 'rn:bridg:roles:rle_1:users', true, done);
      });
    });

    describe('context "/search"', () => {
      it('is allowed access to "/crm-txn/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:crm-txn:customer-profile:_search', true, done);
      });

      it('is allowed access to "/crm-txn/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:crm-txn:customer-profile:_count', true, done);
      });

      it('is allowed access to "/act_1/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_1:customer-profile:_search', true, done);
      });

      it('is allowed access to "/act_2/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_2:customer-profile:_search', true, done);
      });

      it('is allowed access to "/act_1/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_1:customer-profile:_count', true, done);
      });

      it('is allowed access to "/act_2/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_2:customer-profile:_count', true, done);
      });

      it('is NOT allowed access to "/crm-txn/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:crm-txn:_search', false, done);
      });

      it('is NOT allowed access to "/crm-txn/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:crm-txn:_count', false, done);
      });

      it('is NOT allowed access to "/act_1/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_1:_search', false, done);
      });

      it('is NOT allowed access to "/act_2/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_2:_search', false, done);
      });

      it('is NOT allowed access to "/act_1/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_1:_count', false, done);
      });

      it('is NOT allowed access to "/act_2/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_2:_count', false, done);
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:sites', true, done);
      });

      it('is allowed access to "/ste_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', true, done);
      });

      it('is allowed access to "/ste_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', true, done);
      });
    });

    describe('context "/snapshot-fb-export-facebook-accounts"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:snapshot-fb-export-facebook-accounts', true, done);
      });
    });

    describe('context "/snapshot-fb-exports"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:snapshot-fb-exports', true, done);
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:users', true, done);
      });

      it('is allowed access to "/usr_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1', true, done);
      });

      it('is allowed access to "/usr_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2', true, done);
      });

      it('is allowed access to "/usr_1/accounts"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1:accounts', true, done);
      });

      it('is allowed access to "/usr_2/accounts"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2:accounts', true, done);
      });

      it('is allowed access to "/usr_1/authorizations"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1:authorizations', true, done);
      });

      it('is allowed access to "/usr_2/authorizations"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2:authorizations', true, done);
      });

      it('is allowed access to "/usr_1/brands"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1:brands', true, done);
      });

      it('is allowed access to "/usr_2/brands"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2:brands', true, done);
      });

      it('is allowed access to "/usr_1/roles"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1:roles', true, done);
      });

      it('is allowed access to "/usr_2/roles"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2:roles', true, done);
      });
    });

  });

  describe('HTTP "POST"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts', true, done);
      });

      it('is allowed access to "/act_1/actions/add-user"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:actions:add-user', true, done);
      });

      it('is allowed access to "/act_2/actions/add-user"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:actions:add-user', true, done);
      });

      it('is allowed access to "/act_1/brands"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands', true, done);
      });

      it('is allowed access to "/act_2/brands"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands', true, done);
      });

      it('is allowed access to "/act_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:search:customer-profile:_search', true, done);
      });

      it('is allowed access to "/act_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:search:customer-profile:_count', true, done);
      });

      it('is allowed access to "/act_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:search:customer-profile:_search', true, done);
      });

      it('is allowed access to "/act_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:search:customer-profile:_count', true, done);
      });

      it('is allowed access to "/act_1/users/usr_1"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:users:usr_1', true, done);
      });

      it('is allowed access to "/act_2/users/usr_1"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:users:usr_1', true, done);
      });
    });

    describe('context "/analytics"', () => {
      it('is allowed access to "/campaigns/revenue"', (done) => {
        testWardenResponse('create', 'rn:bridg:analytics:campaigns:revenue', true, done);
      });
    });

    describe('context "/authenticate"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:authenticate', true, done);
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is allowed access to "/brd_1/analytics/campaigns/bychannel/facebook"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:analytics:campaigns:bychannel:facebook', true, done);
      });

      it('is allowed access to "/brd_2/analytics/campaigns/bychannel/facebook"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:analytics:campaigns:bychannel:facebook', true, done);
      });

      it('is allowed access to "/brd_1/audiences"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', true, done);
      });

      it('is allowed access to "/brd_2/audiences', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', true, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots', true, done);
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshots', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots', true, done);
      });

      it('is allowed access to "/brd_1/audience-export-groups"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:audience-export-groups', true, done);
      });

      it('is allowed access to "/brd_2/audience-export-groups', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:audience-export-groups', true, done);
      });

      it('is allowed access to "/brd_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_search', true, done);
      });

      it('is allowed access to "/brd_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_count', true, done);
      });

      it('is allowed access to "/brd_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_search', true, done);
      });

      it('is allowed access to "/brd_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_count', true, done);
      });

      it('is allowed access to "/brd_1/snapshot-fb-export-facebook-account"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:snapshot-fb-export-facebook-account', true, done);
      });

      it('is allowed access to "/brd_2/snapshot-fb-export-facebook-account"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:snapshot-fb-export-facebook-account', true, done);
      });

      it('is allowed access to "/brd_1/sites"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:sites', true, done);
      });

      it('is allowed access to "/brd_2/sites', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:sites', true, done);
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:hydra:clients', false, done);
      });
    });

    describe('context "/integrations"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:integrations', true, done);
      });

      it('is allowed access to "/int_1/actions/regenerate-access-key"', (done) => {
        testWardenResponse('create', 'rn:bridg:integrations:int_1:actions:regenerate-access-key', true, done);
      });
    });

    describe('context "/integration-types"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:integration-types', true, done);
      });
    });

    describe('context "/metrics"', () => {
      it('is allowed access to "/write"', (done) => {
        testWardenResponse('create', 'rn:bridg:metrics:write', true, done);
      });

      it('is allowed access to "/query"', (done) => {
        testWardenResponse('create', 'rn:bridg:metrics:query', true, done);
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:hydra:policies', false, done);
      });
    });

    describe('context "/reveal-jobs"', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:reveal-jobs', true, done);
      });
    });

    describe('context "/roles" (rle_1: bridg-admin, rle_2: account-admin, usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles', true, done);
      });

      it('is allowed access to "/rle_1/users/usr_1"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles:rle_1:users:usr_1', true, done);
      });

      it('is allowed access to "/rle_1/users/usr_2"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles:rle_1:users:usr_2', true, done);
      });

      it('is allowed access to "/rle_2/users/usr_1"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles:rle_2:users:usr_1', true, done);
      });

      it('is allowed access to "/rle_2/users/usr_2"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles:rle_2:users:usr_2', true, done);
      });
    });

    describe('context "/scheduler"', () => {
      it('is allowed access to "/run"', (done) => {
        testWardenResponse('create', 'rn:bridg:scheduler:run', true, done);
      });

      it('is allowed access to "/schedule"', (done) => {
        testWardenResponse('create', 'rn:bridg:scheduler:schedule', true, done);
      });
    });

    describe('context "/search"', () => {
      it('is allowed access to "/crm-txn/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:crm-txn:customer-profile:_search', true, done);
      });

      it('is allowed access to "/crm-txn/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:crm-txn:customer-profile:_count', true, done);
      });

      it('is allowed access to "/act_1/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:customer-profile:_search', true, done);
      });

      it('is allowed access to "/act_2/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_2:customer-profile:_search', true, done);
      });

      it('is allowed access to "/act_1/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:customer-profile:_count', true, done);
      });

      it('is allowed access to "/act_2/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_2:customer-profile:_count', true, done);
      });

      it('is NOT allowed access to "/crm-txn/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:crm-txn:_search', false, done);
      });

      it('is NOT allowed access to "/crm-txn/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:crm-txn:_count', false, done);
      });

      it('is NOT allowed access to "/act_1/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:_search', false, done);
      });

      it('is NOT allowed access to "/act_2/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:_search', false, done);
      });

      it('is NOT allowed access to "/act_1/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:_count', false, done);
      });

      it('is NOT allowed access to "/act_2/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_2:_count', false, done);
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:sites', true, done);
      });

      it('is allowed access to "/ste_1/location-links"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1:location-links', true, done);
      });

      it('is allowed access to "/ste_2/location-links"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2:location-links', true, done);
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:users', true, done);
      });

      it('is allowed access to "/usr_1/actions/confirm-account-activation"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_1:actions:confirm-account-activation', true, done);
      });

      it('is allowed access to "/usr_2/actions/confirm-account-activation', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_2:actions:confirm-account-activation', true, done);
      });

      it('is allowed access to "/usr_1/actions/confirm-password-reset"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_1:actions:confirm-password-reset', true, done);
      });

      it('is allowed access to "/usr_2/actions/confirm-password-reset"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_2:actions:confirm-password-reset', true, done);
      });

      it('is allowed access to "/usr_1/actions/request-password-reset"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_1:actions:request-password-reset', true, done);
      });

      it('is allowed access to "/usr_2/actions/request-password-reset"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_2:actions:request-password-reset', true, done);
      });
    });

    describe('context "/warden"', () => {
      it('is NOT allowed access to "/allowed"', (done) => {
        testWardenResponse('create', 'rn:hydra:warden:allowed', false, done);
      });

      it('is NOT allowed access to "/token/allowed"', (done) => {
        testWardenResponse('create', 'rn:hydra:warden:token:allowed', false, done);
      });
    });

  });

  describe('HTTP "PUT"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/act_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1', true, done);
      });

      it('is allowed access to "/act_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2', true, done);
      });

      it('is allowed access to "/act_1/actions/edit-user"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:actions:edit-user', true, done);
      });

      it('is allowed access to "/act_2/actions/edit-user"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:actions:edit-user', true, done);
      });

      it('is NOT allowed access to "/act_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:search:customer-profile:_search', false, done);
      });

      it('is NOT allowed access to "/act_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:search:customer-profile:_count', false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:search:customer-profile:_search', false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:search:customer-profile:_count', false, done);
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is allowed access to "/brd_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1', true, done);
      });

      it('is allowed access to "/brd_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2', true, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1', true, done);
      });

      it('is allowed access to "/brd_2/audiences/aud_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2', true, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports/exp_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports:exp_1', true, done);
      });

      it('is allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports/exp_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports:exp_2', true, done);
      });

      it('is allowed access to "/brd_1/audience-export-groups/aeg_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:audience-export-groups:aeg_1', true, done);
      });

      it('is allowed access to "/brd_2/audience-export-groups/aeg_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:audience-export-groups:aeg_2', true, done);
      });

      it('is NOT allowed access to "/brd_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_search', false, done);
      });

      it('is NOT allowed access to "/brd_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_count', false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_search', false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_count', false, done);
      });
    });

    describe('context "/integration-types"', () => {
      it('is allowed access to "/it_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:integration-types:it_1', true, done);
      });
    });

    describe('context "/roles"', () => {
      it('is allowed access to "/rle_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:roles:rle_1', true, done);
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is allowed access to "/ste_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', true, done);
      });

      it('is allowed access to "/ste_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', true, done);
      });

      it('is allowed access to "/ste_1/location-links/ll_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1:location-links:ll_1', true, done);
      });

      it('is allowed access to "/ste_2/location-links/ll_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_1:location-links:ll_2', true, done);
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/usr_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_1', true, done);
      });

      it('is allowed access to "/usr_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_2', true, done);
      });

      it('is allowed access to "/usr_1/change-email"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_1:change-email', true, done);
      });

      it('is allowed access to "/usr_2/change-email"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_2:change-email', true, done);
      });

      it('is allowed access to "/usr_1/change-password"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_1:change-password', true, done);
      });

      it('is allowed access to "/usr_2/change-password"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_2:change-password', true, done);
      });
    });

  });

  describe('HTTP "DELETE"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/act_1/actions/remove-user"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:actions:remove-user', true, done);
      });

      it('is allowed access to "/act_2/actions/remove-user"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:actions:remove-user', true, done);
      });

      it('is NOT allowed access to "/act_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:search:customer-profile:_search', false, done);
      });

      it('is NOT allowed access to "/act_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:search:customer-profile:_count', false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:search:customer-profile:_search', false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:search:customer-profile:_count', false, done);
      });

      it('is allowed access to "/act_1/users/usr_1"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:users:usr_1', true, done);
      });

      it('is allowed access to "/act_2/users/usr_2"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:users:usr_2', true, done);
      });
    });

    describe('context "/brands"', () => {
      it('is NOT allowed access to "/brd_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_search', false, done);
      });

      it('is NOT allowed access to "/brd_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_count', false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_search', false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_count', false, done);
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/clt_1"', (done) => {
        testWardenResponse('delete', 'rn:hydra:clients:clt_1', false, done);
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/pol_1"', (done) => {
        testWardenResponse('delete', 'rn:hydra:policies:pol_1', false, done);
      });
    });

    describe('context "/roles" (rle_1: bridg-admin, rle_2: account-admin, usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/rle_1/users/usr_1"', (done) => {
        testWardenResponse('delete', 'rn:bridg:roles:rle_1:users:usr_1', true, done);
      });

      it('is allowed access to "/rle_1/users/usr_2"', (done) => {
        testWardenResponse('delete', 'rn:bridg:roles:rle_1:users:usr_2', true, done);
      });

      it('is allowed access to "/rle_2/users/usr_1"', (done) => {
        testWardenResponse('delete', 'rn:bridg:roles:rle_2:users:usr_1', true, done);
      });

      it('is allowed access to "/rle_2/users/usr_2"', (done) => {
        testWardenResponse('delete', 'rn:bridg:roles:rle_2:users:usr_2', true, done);
      });
    });

  });

});
