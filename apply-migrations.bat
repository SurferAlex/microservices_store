@echo off
echo ========================================
echo Применение миграций auth_service
echo ========================================
docker exec -i microservice_postgres psql -U postgres -d auth_service < auth_service\migrations\000001_create_core_table.up.sql
docker exec -i microservice_postgres psql -U postgres -d auth_service < auth_service\migrations\000002_create_ralation_table.up.sql
docker exec -i microservice_postgres psql -U postgres -d auth_service < auth_service\migrations\000003_create_indexes.up.sql
docker exec -i microservice_postgres psql -U postgres -d auth_service < auth_service\migrations\000004_seed_rbac.up.sql

echo.
echo ========================================
echo Применение миграций profile_service
echo ========================================
docker exec -i microservice_postgres psql -U postgres -d profile_service < profile_service\migrations\000001_create_profile.up.sql
docker exec -i microservice_postgres psql -U postgres -d profile_service < profile_service\migrations\000002_create_profile_adresses.up.sql
docker exec -i microservice_postgres psql -U postgres -d profile_service < profile_service\migrations\000003_create_profile_contacts.up.sql

echo.
echo ========================================
echo Все миграции применены успешно!
echo ========================================
pause