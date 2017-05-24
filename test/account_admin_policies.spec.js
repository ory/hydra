
const chai = require('chai');
const  { expect } = chai;
const helper = require('./helper')

describe('The "account-admin" role (act_1 member)', () => {

  const sub = 'account-admin';
  let action;
  let resourceName;
  let response;

  const c1 = { 'account-ids': [['act_1', 'act_1']] };
  const c2 = { 'account-ids': [['act_1', 'act_2']] };

  const updateWardenResponse = (act, rn, cxt, done) => {
    action = act;
    resourceName = rn;
    helper.makeWardenReq(sub, action, resourceName, cxt, (err, res) => {
      response = res;
      done();
    });
  };

  const testWardenResponse = (act, rn, cxt, exp, done) => {
    updateWardenResponse(act, rn, cxt, () => {
      expect(response.body.allowed).to.equal(exp);
      done();
    });
  };

  describe('HTTP "GET"', () => {

    describe('context "/accounts"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts', null, false, done);
      });

      it('is allowed access to "/act_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1', c1, true, done);
      });

      it('is NOT allowed access to "/act_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2', c2, false, done);
      });

      it('is allowed access to "/act_1/brands"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/brands"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands', c2, false, done);
      });

      it('is allowed access to "/act_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:search:customer-profile:_search', c1, true, done);
      });

      it('is allowed access to "/act_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:search:customer-profile:_count', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:search:customer-profile:_search', c2, false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:search:customer-profile:_count', c2, false, done);
      });

      it('is allowed access to "/act_1/sites"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:sites', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/sites"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:sites', c2, false, done);
      });

      it('is allowed access to "/act_1/roles"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:roles', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/roles"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:roles', c2, false, done);
      });

      it('is allowed access to "/act_1/users"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:users', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/users"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:users', c2, false, done);
      });
    });

    describe('context "/audiences"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:audiences', null, false, done);
      });

      it('is NOT allowed access to "/aud_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:audiences:aud_1', null, false, done);
      });

      it('is NOT allowed access to "/aud_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:audiences:aud_2', null, false, done);
      });
    });

    describe('context "/audience-export-groups"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:audience-export-groups', null, false, done);
      });

      it('is NOT allowed access to "/aeg_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:audience-export-groups:aeg_1', null, false, done);
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts', null, false, done);
      });

      it('is allowed access to "/brd_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/analytics/campaigns"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:analytics:campaigns', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/analytics/campaigns"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:analytics:campaigns', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences/aud_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences/aud_2/snapshots"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots/snp_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots:snp_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences/aud_2/snapshots/snp_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots:snp_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports/exp_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports:exp_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports/exp_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports:exp_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/audience-export-groups"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audience-export-groups', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audience-export-groups"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audience-export-groups', c2, false, done);
      });

      it('is allowed access to "/brd_1/audience-export-groups/aeg_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audience-export-groups:aeg_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audience-export-groups/aeg_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audience-export-groups:aeg_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/client-configuration"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:client-configuration', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/client-configuration"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:client-configuration', c2, false, done);
      });

      it('is allowed access to "/brd_1/insights"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:insights', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/insights"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:insights', c2, false, done);
      });

      it('is NOT allowed access to "/brd_1/reveal-jobs"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs', c1, false, done);
      });

      it('is NOT allowed access to "/brd_2/reveal-jobs"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs', c2, false, done);
      });

      it('is NOT allowed access to "/brd_1/reveal-jobs/1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:1', c1, false, done);
      });

      it('is NOT allowed access to "/brd_2/reveal-jobs/2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:2', c2, false, done);
      });

      it('is NOT allowed access to "/brd_1/reveal-jobs/latest"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:latest', c1, false, done);
      });

      it('is NOT allowed access to "/brd_2/reveal-jobs/latest"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:latest', c2, false, done);
      });

      it('is allowed access to "/brd_1/reveal-jobs/latest/artifact"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:reveal-jobs:latest:artifact', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/reveal-jobs/latest/artifact"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:reveal-jobs:latest:artifact', c2, false, done);
      });

      it('is allowed access to "/brd_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_search', c1, true, done);
      });

      it('is allowed access to "/brd_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_count', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_search', c2, false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_count', c2, false, done);
      });

      it('is allowed access to "/brd_1/sites"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/sites"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites', c2, false, done);
      });

      it('is allowed access to "/brd_1/sites/ste_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/sites/ste_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/snapshot-fb-export-facebook-account"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/snapshot-fb-export-facebook-account"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', c2, false, done);
      });
    });

    describe('context "/campaigns"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:campaigns', null, false, done);
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/clt_1"', (done) => {
        testWardenResponse('read', 'rn:hydra:clients:clt_1', null, false, done);
      });
    });

    describe('context "/integrations"', () => {
      it('is NOT allowed access to "/int_1/sync-agent-instances"', (done) => {
        testWardenResponse('read', 'rn:bridg:integrations:int_1:sync-agent-instances', null, false, done);
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:hydra:policies', null, false, done);
      });

      it('is NOT allowed access to "/pol_1"', (done) => {
        testWardenResponse('read', 'rn:hydra:policies:pol_1', null, false, done);
      });
    });

    describe('context "/reveal-jobs"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:reveal-jobs', null, false, done);
      });

      it('is NOT allowed access to "/rvl_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:reveal-jobs:rvl_1', null, false, done);
      });

      it('is NOT allowed access to "/rvl_1/artifact"', (done) => {
        testWardenResponse('read', 'rn:bridg:reveal-jobs:rvl_1:artifact', null, false, done);
      });
    });

    describe('context "/roles"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:roles', null, false, done);
      });

      it('is NOT allowed access to "/rle_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:roles:rle_1', null, false, done);
      });

      it('is NOT allowed access to "/rle_1/users"', (done) => {
        testWardenResponse('read', 'rn:bridg:roles:rle_1:users', null, false, done);
      });
    });

    describe('context "/search"', () => {
      it('is NOT allowed access to "/crm-txn/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:crm-txn:customer-profile:_search', null, false, done);
      });

      it('is NOT allowed access to "/crm-txn/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:crm-txn:customer-profile:_count', null, false, done);
      });

      it('is allowed access to "/act_1/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_1:customer-profile:_search', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/customer-profile/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_2:customer-profile:_search', c2, false, done);
      });

      it('is allowed access to "/act_1/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_1:customer-profile:_count', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/customer-profile/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_2:customer-profile:_count', c2, false, done);
      });

      it('is NOT allowed access to "/crm-txn/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:crm-txn:_search', null, false, done);
      });

      it('is NOT allowed access to "/crm-txn/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:crm-txn:_count', null, false, done);
      });

      it('is NOT allowed access to "/act_1/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_1:_search', c1, false, done);
      });

      it('is NOT allowed access to "/act_2/_search"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_2:_search', c2, false, done);
      });

      it('is NOT allowed access to "/act_1/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_1:_count', c1, false, done);
      });

      it('is NOT allowed access to "/act_2/_count"', (done) => {
        testWardenResponse('read', 'rn:bridg:search:act_2:_count', c2, false, done);
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:sites', null, false, done);
      });

      it('is allowed access to "/ste_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', c1, true, done);
      });

      it('is NOT allowed access to "/ste_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', c2, false, done);
      });
    });

    describe('context "/snapshot-fb-export-facebook-accounts"', () => {
      it('is allowed NOT access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:snapshot-fb-export-facebook-accounts', null, false, done);
      });
    });

    describe('context "/snapshot-fb-exports"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:snapshot-fb-exports', null, false, done);
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('read', 'rn:bridg:users', null, false, done);
      });

      it('is allowed access to "/usr_1"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2', c2, false, done);
      });

      it('is allowed access to "/usr_1/accounts"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1:accounts', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/accounts"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2:accounts', c2, false, done);
      });

      it('is allowed access to "/usr_1/authorizations"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1:authorizations', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/authorizations"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2:authorizations', c2, false, done);
      });

      it('is allowed access to "/usr_1/brands"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1:brands', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/brands"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2:brands', c2, false, done);
      });

      it('is allowed access to "/usr_1/roles"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_1:roles', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/roles"', (done) => {
        testWardenResponse('read', 'rn:bridg:users:usr_2:roles', c2, false, done);
      });
    });

  });

  describe('HTTP "POST"', () => {

    describe('context "/accounts"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts', null, false, done);
      });

      it('is allowed access to "/act_1/actions/add-user"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:actions:add-user', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/actions/add-user"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:actions:add-user', c2, false, done);
      });

      it('is NOT allowed access to "/act_1/brands"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands', c1, false, done);
      });

      it('is NOT allowed access to "/act_2/brands"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands', c2, false, done);
      });

      it('is allowed access to "/act_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:search:customer-profile:_search', c1, true, done);
      });

      it('is allowed access to "/act_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:search:customer-profile:_count', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:search:customer-profile:_search', c2, false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:search:customer-profile:_count', c2, false, done);
      });

      it('is allowed access to "/act_1/users/usr_1"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:users:usr_1', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/users/usr_1"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:users:usr_1', c2, false, done);
      });
    });

    describe('context "/analytics"', () => {
      it('is NOT allowed access to "/campaigns/revenue"', (done) => {
        testWardenResponse('create', 'rn:bridg:analytics:campaigns:revenue', null, false, done);
      });
    });

    describe('context "/authenticate"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:authenticate', null, false, done);
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is allowed access to "/brd_1/analytics/campaigns/bychannel/facebook"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:analytics:campaigns:bychannel:facebook', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/analytics/campaigns/bychannel/facebook"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:analytics:campaigns:bychannel:facebook', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:audiences', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshots"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences/aud_2/snapshots', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshots', c2, false, done);
      });

      it('is allowed access to "/brd_1/audience-export-groups"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:audience-export-groups', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audience-export-groups"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:audience-export-groups', c2, false, done);
      });

      it('is allowed access to "/brd_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_search', c1, true, done);
      });

      it('is allowed access to "/brd_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_count', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_search', c2, false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_count', c2, false, done);
      });

      it('is allowed access to "/brd_1/snapshot-fb-export-facebook-account"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:snapshot-fb-export-facebook-account', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/snapshot-fb-export-facebook-account"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:snapshot-fb-export-facebook-account', c2, false, done);
      });

      it('is allowed access to "/brd_1/sites"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:sites', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/sites', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:sites', c2, false, done);
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:hydra:clients', null, false, done);
      });
    });

    describe('context "/integrations"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:integrations', null, false, done);
      });

      it('is NOT allowed access to "/int_1/actions/regenerate-access-key"', (done) => {
        testWardenResponse('create', 'rn:bridg:integrations:int_1:actions:regenerate-access-key', null, false, done);
      });
    });

    describe('context "/integration-types"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:integration-types', null, false, done);
      });
    });

    describe('context "/metrics"', () => {
      it('is NOT allowed access to "/write"', (done) => {
        testWardenResponse('create', 'rn:bridg:metrics:write', null, false, done);
      });

      it('is allowed access to "/query"', (done) => {
        testWardenResponse('create', 'rn:bridg:metrics:query', null, true, done);
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:hydra:policies', null, false, done);
      });
    });

    describe('context "/reveal-jobs"', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:reveal-jobs', null, false, done);
      });
    });

    describe('context "/roles" (rle_1: bridg-admin, rle_2: account-admin, usr_1: act_1, usr_2: act_2)', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles', null, false, done);
      });

      it('is NOT allowed access to "/rle_1/users/usr_1"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles:499165f9-0e78-42d8-ba72-feb8ee96655d:users:499165f9-0e78-42d8-ba72-feb8ee966666', c1, false, done);
      });

      it('is NOT allowed access to "/rle_1/users/usr_2"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles:499165f9-0e78-42d8-ba72-feb8ee96655d:users:122165f9-0e78-42d8-ba72-feb8ee966666', c2, false, done);
      });

      it('is allowed access to "/rle_2/users/usr_1"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles:999165f9-0e78-42d8-ba72-feb8ee96655d:users:499165f9-0e78-42d8-ba72-feb8ee966666', c1, true, done);
      });

      it('is NOT allowed access to "/rle_2/users/usr_2"', (done) => {
        testWardenResponse('create', 'rn:bridg:roles:999165f9-0e78-42d8-ba72-feb8ee96655d:users:122165f9-0e78-42d8-ba72-feb8ee966666', c2, false, done);
      });
    });

    describe('context "/scheduler"', () => {
      it('is NOT allowed access to "/run"', (done) => {
        testWardenResponse('create', 'rn:bridg:scheduler:run', null, false, done);
      });

      it('is NOT allowed access to "/schedule"', (done) => {
        testWardenResponse('create', 'rn:bridg:scheduler:schedule', null, false, done);
      });
    });

    describe('context "/search"', () => {
      it('is NOT allowed access to "/crm-txn/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:crm-txn:customer-profile:_search', null, false, done);
      });

      it('is NOT allowed access to "/crm-txn/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:crm-txn:customer-profile:_count', null, false, done);
      });

      it('is allowed access to "/act_1/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:customer-profile:_search', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/customer-profile/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_2:customer-profile:_search', c2, false, done);
      });

      it('is allowed access to "/act_1/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:customer-profile:_count', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/customer-profile/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_2:customer-profile:_count', c2, false, done);
      });

      it('is NOT allowed access to "/crm-txn/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:crm-txn:_search', null, false, done);
      });

      it('is NOT allowed access to "/crm-txn/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:crm-txn:_count', null, false, done);
      });

      it('is NOT allowed access to "/act_1/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:_search', c1, false, done);
      });

      it('is NOT allowed access to "/act_2/_search"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_2:_search', c2, false, done);
      });

      it('is NOT allowed access to "/act_1/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_1:_count', c1, false, done);
      });

      it('is NOT allowed access to "/act_2/_count"', (done) => {
        testWardenResponse('create', 'rn:bridg:search:act_2:_count', c2, false, done);
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:sites', null, false, done);
      });

      it('is allowed access to "/ste_1/location-links"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1:location-links', c1, true, done);
      });

      it('is NOT allowed access to "/ste_2/location-links"', (done) => {
        testWardenResponse('create', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2:location-links', c2, false, done);
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is NOT allowed access to "/"', (done) => {
        testWardenResponse('create', 'rn:bridg:users', null, false, done);
      });

      it('is allowed access to "/usr_1/actions/confirm-account-activation"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_1:actions:confirm-account-activation', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/actions/confirm-account-activation', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_2:actions:confirm-account-activation', c2, false, done);
      });

      it('is allowed access to "/usr_1/actions/confirm-password-reset"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_1:actions:confirm-password-reset', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/actions/confirm-password-reset"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_2:actions:confirm-password-reset', c2, false, done);
      });

      it('is allowed access to "/usr_1/actions/request-password-reset"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_1:actions:request-password-reset', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/actions/request-password-reset"', (done) => {
        testWardenResponse('create', 'rn:bridg:users:usr_2:actions:request-password-reset', c2, false, done);
      });
    });

    describe('context "/warden"', () => {
      it('is NOT allowed access to "/allowed"', (done) => {
        testWardenResponse('create', 'rn:hydra:warden:allowed', null, false, done);
      });

      it('is NOT allowed access to "/token/allowed"', (done) => {
        testWardenResponse('create', 'rn:hydra:warden:token:allowed', null, false, done);
      });
    });

  });

  describe('HTTP "PUT"', () => {

    describe('context "/accounts"', () => {
      it('is NOT allowed access to "/act_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1', c1, false, done);
      });

      it('is NOT allowed access to "/act_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2', c2, false, done);
      });

      it('is allowed access to "/act_1/actions/edit-user"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:actions:edit-user', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/actions/edit-user"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:actions:edit-user', c2, false, done);
      });

      it('is NOT allowed access to "/act_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:search:customer-profile:_search', c1, false, done);
      });

      it('is NOT allowed access to "/act_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:search:customer-profile:_count', c1, false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:search:customer-profile:_search', c2, false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:search:customer-profile:_count', c2, false, done);
      });
    });

    describe('context "/brands" (brd_1: act_1, brd_2: act_2)', () => {
      it('is NOT allowed access to "/brd_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1', c1, false, done);
      });

      it('is NOT allowed access to "/brd_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences/aud_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/audience-export-groups/aeg_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:audience-export-groups:aeg_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audience-export-groups/aeg_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:audience-export-groupss:aud_2', c2, false, done);
      });

      it('is allowed access to "/brd_1/audiences/aud_1/snapshot-fb-exports/exp_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports:exp_1', c1, true, done);
      });

      it('is NOT allowed access to "/brd_2/audiences/aud_2/snapshot-fb-exports/exp_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:audiences:aud_2:snapshot-fb-exports:exp_2', c2, false, done);
      });

      it('is NOT allowed access to "/brd_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_search', c1, false, done);
      });

      it('is NOT allowed access to "/brd_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_count', c1, false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_search', c2, false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_count', c2, false, done);
      });
    });

    describe('context "/integration-types"', () => {
      it('is NOT allowed access to "/it_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:integration-types:it_1', null, false, done);
      });
    });

    describe('context "/roles"', () => {
      it('is NOT allowed access to "/rle_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:roles:rle_1', null, false, done);
      });
    });

    describe('context "/sites" (ste_1: act_1 - brd_1, ste_2: act_2 - brd_2)', () => {
      it('is allowed access to "/ste_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', c1, true, done);
      });

      it('is NOT allowed access to "/ste_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_2', c2, false, done);
      });

      it('is allowed access to "/ste_1/location-links/ll_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1:location-links:ll_1', c1, true, done);
      });

      it('is NOT allowed access to "/ste_2/location-links/ll_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:accounts:act_2:brands:brd_2:sites:ste_1:location-links:ll_2', c2, false, done);
      });
    });

    describe('context "/users" (usr_1: act_1, usr_2: act_2)', () => {
      it('is allowed access to "/usr_1"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_1', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_2', c2, false, done);
      });

      it('is allowed access to "/usr_1/change-email"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_1:change-email', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/change-email"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_2:change-email', c2, false, done);
      });

      it('is allowed access to "/usr_1/change-password"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_1:change-password', c1, true, done);
      });

      it('is NOT allowed access to "/usr_2/change-password"', (done) => {
        testWardenResponse('update', 'rn:bridg:users:usr_2:change-password', c2, false, done);
      });
    });

  });

  describe('HTTP "DELETE"', () => {

    describe('context "/accounts"', () => {
      it('is allowed access to "/act_1/actions/remove-user"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:actions:remove-user', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/actions/remove-user"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:actions:remove-user', c2, false, done);
      });

      it('is NOT allowed access to "/act_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:search:customer-profile:_search', c1, false, done);
      });

      it('is NOT allowed access to "/act_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:search:customer-profile:_count', c1, false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:search:customer-profile:_search', c2, false, done);
      });

      it('is NOT allowed access to "/act_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:search:customer-profile:_count', c2, false, done);
      });

      it('is allowed access to "/act_1/users/usr_1"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:users:usr_1', c1, true, done);
      });

      it('is NOT allowed access to "/act_2/users/usr_2"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:users:usr_2', c2, false, done);
      });
    });

    describe('context "/brands"', () => {
      it('is NOT allowed access to "/brd_1/search/customer-profile/_search"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_search', c1, false, done);
      });

      it('is NOT allowed access to "/brd_1/search/customer-profile/_count"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_1:brands:brd_1:search:customer-profile:_count', c1, false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_search"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_search', c2, false, done);
      });

      it('is NOT allowed access to "/brd_2/search/customer-profile/_count"', (done) => {
        testWardenResponse('delete', 'rn:bridg:accounts:act_2:brands:brd_2:search:customer-profile:_count', c2, false, done);
      });
    });

    describe('context "/clients"', () => {
      it('is NOT allowed access to "/clt_1"', (done) => {
        testWardenResponse('delete', 'rn:hydra:clients:clt_1', null, false, done);
      });
    });

    describe('context "/policies"', () => {
      it('is NOT allowed access to "/pol_1"', (done) => {
        testWardenResponse('delete', 'rn:hydra:policies:pol_1', null, false, done);
      });
    });

    describe('context "/roles" (rle_1: bridg-admin, rle_2: account-admin, usr_1: act_1, usr_2: act_2)', () => {
      it('is NOT allowed access to "/rle_1/users/usr_1"', (done) => {
        testWardenResponse('delete', 'rn:bridg:roles:499165f9-0e78-42d8-ba72-feb8ee96655d:users:499165f9-0e78-42d8-ba72-feb8ee966666', c1, false, done);
      });

      it('is NOT allowed access to "/rle_1/users/usr_2"', (done) => {
        testWardenResponse('delete', 'rn:bridg:roles:499165f9-0e78-42d8-ba72-feb8ee96655d:users:122165f9-0e78-42d8-ba72-feb8ee966666', c2, false, done);
      });

      it('is allowed access to "/rle_2/users/usr_1"', (done) => {
        testWardenResponse('delete', 'rn:bridg:roles:999165f9-0e78-42d8-ba72-feb8ee96655d:users:499165f9-0e78-42d8-ba72-feb8ee966666', c1, true, done);
      });

      it('is NOT allowed access to "/rle_2/users/usr_2"', (done) => {
        testWardenResponse('delete', 'rn:bridg:roles:999165f9-0e78-42d8-ba72-feb8ee96655d:users:122165f9-0e78-42d8-ba72-feb8ee966666', c2, false, done);
      });
    });

  });

});
