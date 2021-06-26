CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


CREATE TABLE public.models (
    model_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    model_name character varying(255) NOT NULL,
    windfarm_id uuid NOT NULL,
    value double precision DEFAULT 0 NOT NULL,
    icuf double precision DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE public.outputs (
    output_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    turbine_id uuid NOT NULL,
    speed double precision NOT NULL,
    production double precision NOT NULL
);

CREATE TABLE public.productions (
    production_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    date date,
    "time" time(6) without time zone,
    icuf double precision NOT NULL,
    value double precision,
    turbines_models_id uuid NOT NULL,
    wind_speed double precision DEFAULT 0 NOT NULL,
    altitude double precision NOT NULL,
    wind_direction double precision NOT NULL,
    shading double precision NOT NULL,
    speed_with_shading double precision NOT NULL
);


ALTER TABLE public.productions OWNER TO psihachina;

CREATE TABLE public.times (
    "time" time(6) without time zone
);

CREATE TABLE public.turbines (
    turbine_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    turbine_name character varying(255) NOT NULL,
    maximum_power double precision NOT NULL,
    max_wind_speed double precision NOT NULL,
    min_wind_speed double precision NOT NULL,
    tower_height double precision NOT NULL,
    number_blades double precision NOT NULL,
    rotor_diameter double precision NOT NULL
);

CREATE TABLE public.turbines_models (
    turbines_models_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    turbine_name character varying,
    latitude double precision,
    longitude double precision,
    model_id uuid NOT NULL,
    y double precision NOT NULL,
    x double precision NOT NULL,
    z double precision NOT NULL,
    value double precision DEFAULT 0 NOT NULL,
    icuf double precision DEFAULT 0 NOT NULL
);

CREATE TABLE public.users (
    user_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    registered_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    admin_confirm boolean DEFAULT false NOT NULL,
    email_confirm boolean DEFAULT false NOT NULL
);

CREATE TABLE public.users_windfarms (
    users_windfarms_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    windfarm_id uuid NOT NULL
);

CREATE TABLE public.windfarms (
    windfarm_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    windfarm_name character varying(255) NOT NULL,
    polygon polygon,
    longitude double precision NOT NULL,
    latitude double precision NOT NULL,
    capacity double precision NOT NULL,
    description character varying(255) NOT NULL,
    polygon_radius double precision NOT NULL,
    altitude double precision NOT NULL
);

CREATE TABLE public.winds (
    wind_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    windfarm_id uuid NOT NULL,
    date date NOT NULL,
    "time" time(6) without time zone NOT NULL,
    wind_speed double precision NOT NULL,
    wind_direction double precision NOT NULL,
    altitude double precision NOT NULL
);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT email_unique UNIQUE (email);

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (model_id);

ALTER TABLE ONLY public.outputs
    ADD CONSTRAINT outputs_pkey PRIMARY KEY (output_id);

ALTER TABLE ONLY public.productions
    ADD CONSTRAINT productions_pkey PRIMARY KEY (production_id);

ALTER TABLE ONLY public.turbines_models
    ADD CONSTRAINT turbines_models_pkey PRIMARY KEY (turbines_models_id);

ALTER TABLE ONLY public.turbines
    ADD CONSTRAINT turbines_pkey PRIMARY KEY (turbine_id);

ALTER TABLE ONLY public.times
    ADD CONSTRAINT unique_time UNIQUE ("time");

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);

ALTER TABLE ONLY public.users_windfarms
    ADD CONSTRAINT users_windfarms_pkey PRIMARY KEY (users_windfarms_id);

ALTER TABLE ONLY public.windfarms
    ADD CONSTRAINT windfarms_pkey PRIMARY KEY (windfarm_id);

ALTER TABLE ONLY public.winds
    ADD CONSTRAINT winds_pkey PRIMARY KEY (wind_id);

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_windfarm_id_fkey FOREIGN KEY (windfarm_id) REFERENCES public.windfarms(windfarm_id) ON DELETE CASCADE;


ALTER TABLE ONLY public.outputs
    ADD CONSTRAINT outputs_turbine_id_fkey FOREIGN KEY (turbine_id) REFERENCES public.turbines(turbine_id) ON DELETE CASCADE;

ALTER TABLE ONLY public.productions
    ADD CONSTRAINT productions_turbines_models_id_fkey FOREIGN KEY (turbines_models_id) REFERENCES public.turbines_models(turbines_models_id) ON DELETE CASCADE;


ALTER TABLE ONLY public.turbines_models
    ADD CONSTRAINT turbines_models_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(model_id) ON DELETE CASCADE;

ALTER TABLE ONLY public.turbines
    ADD CONSTRAINT turbines_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;

ALTER TABLE ONLY public.users_windfarms
    ADD CONSTRAINT users_windfarms_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;

ALTER TABLE ONLY public.users_windfarms
    ADD CONSTRAINT users_windfarms_windfarm_id_fkey FOREIGN KEY (windfarm_id) REFERENCES public.windfarms(windfarm_id) ON DELETE CASCADE;

ALTER TABLE ONLY public.winds
    ADD CONSTRAINT winds_windfarm_id_fkey FOREIGN KEY (windfarm_id) REFERENCES public.windfarms(windfarm_id) ON DELETE CASCADE;


