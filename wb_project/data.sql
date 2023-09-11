CREATE TABLE
    public.delivery (
        id INTEGER PRIMARY KEY,
        user_id INTEGER unique FOREIGN KEY REFERENCES public.user (id),
        name VARCHAR(100) NOT NULL,
        phone VARCHAR(100) NOT NULL,
        zip VARCHAR(100) NOT NULL,
        city VARCHAR(100) NOT NULL,
        address VARCHAR(100) NOT NULL,
        region VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL,
    );

CREATE TABLE
    public.payment (
        id INTEGER PRIMARY KEY,
        user_id INTEGER unique FOREIGN KEY REFERENCES public.user (id),
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
    );

CREATE TABLE
    public.items (
        id INTEGER PRIMARY KEY,
        chrt_id INTEGER,
        track_number FOREIGN KEY REFERENCES public.user (track_number),
        price INTEGER,
        rid VARCHAR(100) NOT NULL,
        name VARCHAR(100) NOT NULL,
        sale INTEGER,
        size VARCHAR(100) NOT NULL,
        total_price INTEGER,
        nm_id INTEGER,
        brand VARCHAR(100) NOT NULL,
        status INTEGER,
    );

CREATE TABLE
    public.user (
        id INTEGER PRIMARY KEY,
        order_uuid VARCHAR(100) NOT NULL,
        track_number VARCHAR(100) NOT NULL,
        entry VARCHAR(100) NOT NULL,
        locale VARCHAR(100) NOT NULL,
        internal_signature VARCHAR(100) NOT NULL,
        customer_id VARCHAR(100) NOT NULL,
        delivery_service VARCHAR(100) NOT NULL,
        shardkey VARCHAR(100) NOT NULL,
        sm_id INTEGER,
        date_created VARCHAR(50) NOT NULL,
        oof_shard VARCHAR(100) NOT NULL,
        CONSTRAINT user_PK PRIMARY KEY (id)
    );