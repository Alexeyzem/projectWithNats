CREATE TABLE
    public.delivery (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES public.users (id),
        name VARCHAR(100) NOT NULL,
        phone VARCHAR(100) NOT NULL,
        zip VARCHAR(100) NOT NULL,
        city VARCHAR(100) NOT NULL,
        address VARCHAR(100) NOT NULL,
        region VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL,
        UNIQUE (user_id)
    );

CREATE TABLE
    public.payment (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES public.users (id),
        transaction VARCHAR(100) NOT NULL,
        request_id VARCHAR(100) NOT NULL,
        currency VARCHAR(100) NOT NULL,
        provider VARCHAR(100) NOT NULL,
        amount INTEGER,
        payment_dt INTEGER,
        bank VARCHAR(100) NOT NULL,
        delivery_cost INTEGER,
        goods_total INTEGER,
        custom_fee INTEGER,
        UNIQUE (user_id)
    );

CREATE TABLE
    public.items (
        id SERIAL PRIMARY KEY,
        chrt_id INTEGER,
        track_number VARCHAR(100) REFERENCES public.users (track_number),
        price INTEGER,
        rid VARCHAR(100) NOT NULL,
        name VARCHAR(100) NOT NULL,
        sale INTEGER,
        size VARCHAR(100) NOT NULL,
        total_price INTEGER,
        nm_id INTEGER,
        brand VARCHAR(100) NOT NULL,
        status INTEGER
    );

CREATE TABLE
    users (
        id SERIAL PRIMARY KEY,
        order_uuid VARCHAR(19) NOT NULL,
        track_number VARCHAR(100),
        entry VARCHAR(100),
        locale VARCHAR(100),
        internal_signature VARCHAR(100),
        customer_id VARCHAR(100),
        delivery_service VARCHAR(100),
        shardkey VARCHAR(100),
        sm_id INTEGER,
        date_created VARCHAR(50),
        oof_shard VARCHAR(100),
        UNIQUE (track_number)
    );