CREATE TABLE "files" (
  "id" BIGSERIAL PRIMARY KEY,
  "filename" VARCHAR,
  "secondname" VARCHAR,
  "path" VARCHAR,
  "size" INT,
  "uploaded_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);