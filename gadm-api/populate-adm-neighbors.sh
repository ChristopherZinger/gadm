

SERVICE_TYPE="cron_job"
CRON_JOB_NAME="populate_adm_neighbors"

export SERVICE_TYPE
export CRON_JOB_NAME

cd main && go run .
