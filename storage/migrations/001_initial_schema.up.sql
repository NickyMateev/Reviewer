BEGIN;

CREATE TABLE projects (
  id SERIAL PRIMARY KEY NOT NULL,
  name varchar(255) NOT NULL,
  repo_name varchar(255) NOT NULL,
  repo_owner varchar(255) NOT NULL,
  UNIQUE (repo_name, repo_owner)
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY NOT NULL,
  username varchar(255) NOT NULL UNIQUE,
  github_id INTEGER NOT NULL
);

CREATE TABLE project_users (
  project_id INTEGER,
  user_id INTEGER,
  PRIMARY KEY (project_id, user_id),
  CONSTRAINT project_users_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT project_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE pull_requests (
  id SERIAL PRIMARY KEY NOT NULL,
  user_id INTEGER NOT NULL,
  project_id INTEGER NOT NULL,
  title varchar(255) NOT NULL,
  url varchar(255) NOT NULL,
  number INTEGER NOT NULL,
  github_id INTEGER NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT pull_requests_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT pull_requests_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE approvers (
  user_id INTEGER,
  pull_request_id INTEGER,
  PRIMARY KEY (user_id, pull_request_id),
  CONSTRAINT approvers_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT approvers_pull_request_id_fkey FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE
);

CREATE TABLE commenters (
  user_id INTEGER,
  pull_request_id INTEGER,
  PRIMARY KEY (user_id, pull_request_id),
  CONSTRAINT commenters_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT commenters_pull_request_id_fkey FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE
);

CREATE TABLE idlers (
  user_id INTEGER,
  pull_request_id INTEGER,
  PRIMARY KEY (user_id, pull_request_id),
  CONSTRAINT idlers_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT idlers_pull_request_id_fkey FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE
);

CREATE TABLE reviewers (
  user_id INTEGER,
  pull_request_id INTEGER,
  PRIMARY KEY (user_id, pull_request_id),
  CONSTRAINT reviewers_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT reviewers_pull_request_id_fkey FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE
);

COMMIT;
