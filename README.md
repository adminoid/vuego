# Vuego cmf

Content management system with universal tree data structure.
Based on golang, vue.js framework and postgresql database.

# Creating migrations

`migrate create -ext sql -dir db/migrations -seq create_users_table`
then
`make up` or `make down`
