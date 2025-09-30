# db

Raw psql queries.


## Create migration
```bash
 migrate create -ext sql -dir _migrations/ -seq name
```

## Configuring local psql

To configure PostgreSQL for local authentication with a username and password on Ubuntu:

Set a password for the PostgreSQL user:

```
sudo -u postgres psql
\password postgres
```

Enter the desired password.

Edit pg_hba.conf to use password authentication:

Open the file (location: /etc/postgresql/<version>/main/pg_hba.conf).

Change local and/or host lines for your user/database to use md5 or scram-sha-256:

```
local   all   postgres   md5
host    all   postgres   127.0.0.1/32   md5
```
Save and exit.

Reload PostgreSQL:

```
sudo systemctl reload postgresql
```

Now you can connect using:

```
psql -U postgres -h localhost -W
```

You'll be prompted for the password.

## Migrating local database

```bash
go run ./cmd/example migrate
```
