FROM postgres:15


ENV POSTGRES_USER postgres      
ENV POSTGRES_PASSWORD my_password  
ENV POSTGRES_DB fc2            


COPY init.sql /docker-entrypoint-initdb.d/


EXPOSE 5432