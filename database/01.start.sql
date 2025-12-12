DROP TABLE IF EXISTS public.students;

CREATE TABLE public.students
(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "email_lower" VARCHAR(100) GENERATED ALWAYS AS (lower("email")) STORED,
    UNIQUE ("email_lower")
);