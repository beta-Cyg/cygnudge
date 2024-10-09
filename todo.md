- Judge module for server (working)
- Network module for server (working)
- Network module for client (working)

- Enable autosave for redis (new shell script) (partially done)
- Run judge tasks in docker container

- Enable EOF handling for server (chcnage Fatals to Prints)
- If it has loged in, firstly logout and then login

- Move all configurations to {}.json, rather than using macros in headers

- Refactor test using standard library; change into unit testing.

# Bugs
- zip & unzip module based on python (Solution: Use Golang pack module)
- reduce unnecessary directory in zip file created by bin/cygpack.py
