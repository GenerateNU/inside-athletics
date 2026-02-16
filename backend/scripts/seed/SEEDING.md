# Seeding

Here are the available commands for data seeding:

## Main command (Generation and Seeding)

```
make seed
```

Main seeding command, generates the college and sports data and populates the databases (You'll probably just want to run this)

## College Data

```
make generate-colleges-data
```

Generates college data and inputs it into `colleges.json`

## Sports Data

```
make generate-sports-data
```

Generates sport data and inputs it into `sports.json`

## Generation

```
make generate-seed-data
```

Generates both the college and the sports data (basically runs the two commmands above)

---

**Note:** Python dependencies (pandas, kagglehub) are automatically installed when running any of these commands.
