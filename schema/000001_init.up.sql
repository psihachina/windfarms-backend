CREATE EXTENSION "uuid-ossp";

CREATE TABLE users (
    user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);
CREATE TABLE windfarms (
    windfarm_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    windfarm_name VARCHAR(255) NOT NULL,
    polygon POLYGON,
    longitude DOUBLE PRECISION NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    capacity DOUBLE PRECISION NOT NULL,
    range_to_city DOUBLE PRECISION NOT NULL,
    range_to_road DOUBLE PRECISION NOT NULL,
    range_to_city_line DOUBLE PRECISION NOT NULL,
    city_longitude DOUBLE PRECISION NOT NULL,
    city_latitude DOUBLE PRECISION NOT NULL,
    description VARCHAR(255) NOT NULL
);
CREATE TABLE winds(
    wind_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    windfarm_id uuid references windfarms(windfarm_id) on delete cascade not null,
    date DATE NOT NULL,
    time Time NOT NULL,
    temperature DOUBLE PRECISION NOT NULL,
    wind_speed DOUBLE PRECISION NOT NULL,
    wind_direction VARCHAR(255) NOT NULL,
    humidity DOUBLE PRECISION NOT NULL,
    altitude DOUBLE PRECISION NOT NULL
);

CREATE TABLE users_windfarms (
    users_windfarms_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid references users(user_id) on delete cascade not null,
    windfarm_id uuid references windfarms(windfarm_id) on delete cascade not null
);

CREATE TABLE turbines(
    turbine_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid references users(user_id) on delete cascade not null,
    turbine_name VARCHAR(255) NOT NULL,
    maximum_power DOUBLE PRECISION NOT NULL,
    max_wind_speed DOUBLE PRECISION NOT NULL,
    min_wind_speed DOUBLE PRECISION NOT NULL,
    tower_height DOUBLE PRECISION NOT NULL,
    number_blades DOUBLE PRECISION NOT NULL,
    rotor_diameter DOUBLE PRECISION NOT NULL,
    annual_turbine_maintenance DOUBLE PRECISION  NOT NULL
);


CREATE TABLE outputs(
    output_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    turbine_id uuid references turbines(turbine_id) on delete cascade not null,
    speed DOUBLE PRECISION NOT NULL,
    production DOUBLE PRECISION NOT NULL
);