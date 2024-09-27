ServiceSentry

Check if cratedb is still running ok, if not restart the service and report error to log file in /var/log/servicesentry/error.log

Add task to cronjob and check every minute

\* \* \* \* \* cd /yourdir && ./servicesentry

Keep in mind that the service will be started even if you stop the cratedb manually 
