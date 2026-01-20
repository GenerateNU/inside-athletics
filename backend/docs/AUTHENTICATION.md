## How Authentication works and how to authenticate

Every requesto our backend needs to be authenticated. This is essential to making sure our apis are safe and bad actors cannot pretend to be other users. 

We use the token supabase give the user through login for this. I created a helper target called gen-token which will encode one of these tokens giving the user-id. Now I can either manually write out the Authorization header or just paste it into the spot in the documentation. 

To use this value once the Auth Middleware verifies it fetch it from the ctx. 


