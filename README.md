# mwx-go-auth-service


# Tables

```
users
-----
id (UUID)
email
display_name
avatar_url
role
status
created_at
updated_at

user_identities
---------------
id
user_id
provider
provider_user_id
password_hash
created_at
updated_at
```
