# Airflow Version to REST API Version Matrix

## Summary

| Airflow Version | REST API Version | Base Path | Status |
|-----------------|------------------|-----------|--------|
| 1.x (<=1.10.x) | Experimental API | `/api/experimental/` | Deprecated in 2.0 |
| 2.0.x - 2.11.x | Stable API v1 | `/api/v1/` | Current for Airflow 2.x |
| 3.0.x+ | API v2 | `/api/v2/` | Introduced in Airflow 3.0 |

## Detailed Breakdown

### Airflow 1.x
- **API**: Experimental REST API
- **Path**: `/api/experimental/`
- **Notes**: 
  - No stable REST API
  - Experimental API disabled by default in Airflow 2.0+
  - Enable with: `enable_experimental_api = True` in `[api]` section

### Airflow 2.0 - 2.11.x (Current Stable)
- **API**: Stable REST API v1
- **Path**: `/api/v1/`
- **Key Versions**:
  - **2.0.0** (Dec 17, 2020): Introduced stable API v1, deprecated experimental API
  - **2.1.0**: Added new endpoints
  - **2.10.x** (e.g., 2.10.2, 2.10.5): Latest 2.x series with API v1
  - **2.11.x**: Last minor release in 2.x series (Limited maintenance until Apr 2026)

### Airflow 3.0+ (Future/Current)
- **API**: REST API v2
- **Path**: `/api/v2/`
- **Key Changes**:
  - **3.0.0** (Apr 22, 2025): Introduced API v2
  - New endpoints: `/assets`, `/backfills`, `/monitor/health`, etc.
  - Breaking changes from v1 (see [Airflow 2 to 3 REST API Changes](https://github.com/apache/airflow/issues/43378))

## API Client Compatibility

| Component | Airflow 2.x (v1) | Airflow 3.x (v2) |
|-----------|-------------------|-------------------|
| OpenAPI Spec | `v1.yaml` or `v1-rest-api-generated.yaml` | `v2-rest-api-generated.yaml` |
| Base URL | `http://localhost:8080/api/v1/` | `http://localhost:8080/api/v2/` |
| Auth | `airflow.api.auth.backend.basic_auth` | Same + new auth endpoints |

## References

- [Airflow REST API Documentation](https://airflow.apache.org/docs/apache-airflow/stable/stable-rest-api-ref.html)
- [Upgrading to Airflow 2.0+](https://airflow.apache.org/docs/apache-airflow/2.0.0/upgrading-to-2.html)
- [Airflow 2 to 3 REST API Changes](https://github.com/apache/airflow/issues/43378)
- [Google Cloud Composer - Airflow API Versions](https://cloud.google.com/composer/docs/composer-2/access-airflow-api)

## Notes for airflow-tui

- **Current code** (`generated_client.go`) points to `/api/v2/` (Airflow 3.x)
- **Docker setup** uses Airflow 2.10.2 which only has `/api/v1/`
- **Fix needed**: Change URL to `/api/v1/` OR implement client versioning based on `api_version` config
- **OpenAPI spec** in repo (`airflow-v2-openapi.yml`) uses `/api/v1/` paths - this is confusingly named
