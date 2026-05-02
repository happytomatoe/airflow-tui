FROM apache/airflow:2.10.2

ENV AIRFLOW__CORE__EXECUTOR=LocalExecutor
ENV AIRFLOW__DATABASE__SQL_ALCHEMY_CONN=postgresql+psycopg2://airflow:airflow@postgres/airflow
ENV AIRFLOW__CORE__FERNET_KEY=
ENV AIRFLOW__CORE__DAGS_ARE_PAUSED_AT_CREATION=False
ENV AIRFLOW__CORE__LOAD_EXAMPLES=True
ENV AIRFLOW__API__AUTH_BACKENDS=airflow.api.auth.backend.basic_auth

USER root
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
USER airflow