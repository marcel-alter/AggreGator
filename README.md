Project gator!
This is a Blog Aggregator CLI tool made in Go for educational purposes.
It works only with RSS links for now.
for gator to work, we need Go aswell as PostgreSQL to be installed.

for PostgreSQL on Linux use command: 
    sudo apt update
    sudo apt install postgresql postgresql-contrib
to install gator use command:
    go install github.com/marcel-alter/AggreGator@latest

To use gator you will first need to register a User.
Use command: register <UserName>
you can use multiple users.
to switch between registered users, use command: login <UserName>
To view all users, use command: users

To add a Feed to the Aggregator, use command: addfeed <Name> <FeedURL>
to view all Feeds and their owners, use command: feeds
The current user can follow a Feed by using the command: follow <FeedURL>


With multiple feeds being stored in the database, we can aggregate the feeds posts.
Use command: agg <duration>
example: "agg 1h10m5s" to aggregate articles every 1hour 10minutes and 5 seconds.

To view the now aggregated posts, use command: browse
This command displays only the posts of the feed the current user is following.

For a detailed list on all available commands, use command: help