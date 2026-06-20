create table if not exists users (
       id serial primary key,
       is_active boolean not null default true,
       created_at timestamptz not null default now(),
       updated_at timestamptz not null default now(),
       user_id bigint not null,
       first_name text not null,
       last_name text,
       username text,
       score integer not null default 0
);

create or replace function updated_at()
returns trigger as $$
begin
    new.updated_at = current_timestamp;
    return new;
end;
$$ language plpgsql;

drop trigger if exists users_update_at on users;
create trigger users_update_at
       before update on users
       for each row
       execute function updated_at();
