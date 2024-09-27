ServiceSentry

Check if cratedb is still running ok, if not restart the service and report error to log file in /var/log/servicecentry/error.log

Add task to cronjob and check every minute

\* \* \* \* \* cd /yourdir && ./servicesentry
