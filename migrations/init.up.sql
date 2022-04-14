create table sale
(
    id              serial primary key,
    uuid            uuid,
    name            text not null,
    start_date      date not null,
    end_date        date not null,
    stop_factor     smallint,
    contractor_type smallint,
    type            smallint,
    state           smallint,
    created_by      uuid,
    created_at      date
);

create table sale_product
(
    id              serial primary key,
    uuid            uuid,
    sale_id         int, -- uuid ?
    region_id       int,
    nomencl_uuid    uuid,
    price           int,
    max_count       int     default 0,
    max_order_count int     default 0,
    stock_available int     default 0,
    available       int     default 0,
    is_feed         boolean default false,
    status          int     default 0,
    created_by      uuid,
    created_at      date,
    updated_by      uuid,
    updated_at      date
);

create table sell
(
    id              serial primary key,
    sale_product_id bigint,
    order_uuid      uuid,
    quantity        int,
    created_at      date
)

-- todo индексы
-- todo comment on column