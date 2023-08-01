create table user(
    user_id varchar(32),
    name varchar(32),
    score integer,
    energy integer,
    update_time bigint,
    score_time bigint,
    create_time bigint,
    ad_times integer,
    invite_times integer,
    other_times integer,
    invitor varchar(32)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;

create table server(
    reset_time integer
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;
