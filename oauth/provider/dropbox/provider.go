package dropbox

// 			/oauth2/auth?
//				response_type: code
// 				client_id: client_id
// 				scope: ar,Scopes
//				// we're using the same state throughout the process, this should actually not be included in the redirect url, because most oauth providers will return the given state
//				state: ar.State
// 				provider={provider}
// 				redirect_uri={original_redirect)
