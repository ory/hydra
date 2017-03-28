
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

    describe('context "/brands"', () => {
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
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/analytics/campaigns"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:analytics:campaigns', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/analytics/campaigns"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:analytics:campaigns', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/audiences"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/audiences"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/audiences/aud_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/audiences/aud_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:audiences:aud_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/audiences/aud_1/snapshots"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/audiences/aud_2/snapshots"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:audiences:aud_2:snapshots', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/audiences/aud_1/snapshots/snp_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshots:snp_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/audiences/aud_2/snapshots/snp_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:audiences:aud_2:snapshots:snp_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/audiences/aud_1/snapshot-fb-exports"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/audiences/aud_2/snapshot-fb-exports"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:audiences:aud_2:snapshot-fb-exports', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/audiences/aud_1/snapshot-fb-exports/exp_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences:aud_1:snapshot-fb-exports:exp_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/audiences/aud_2/snapshot-fb-exports/exp_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:audiences:aud_2:snapshot-fb-exports:exp_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });


      it('is allowed access to "/brands/brd_1/sites"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/sites"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:sites', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/sites/ste_1"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:sites:ste_1', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/sites/ste_2"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:sites:ste_2', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_1/snapshot-fb-export-facebook-account"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_1:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });

      it('is allowed access to "/brands/brd_2/snapshot-fb-export-facebook-account"', (done) => {
        updateWardenResponse('read', 'rn:bridg:accounts:act_1:brands:brd_2:audiences', () => {
          expect(response.body.allowed).to.equal(true);
          done();
        });
      });
    });

  });

});
