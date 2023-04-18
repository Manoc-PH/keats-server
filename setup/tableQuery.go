// --create table nutrient(id int, name varchar, unit_name varchar, nutrient_nbr float8, rank int);

// --COPY nutrient FROM 'D:\.FOR PRODUCTIVITY\Manoc\FoodData Central\csv\nutrient.csv' WITH CSV HEADER;

// --

// --create table food_nutrient(id int, fdc_id int, nutrient_id int, amount float8, data_points varchar, derivation_id varchar, min varchar, max varchar, median varchar, loq varchar, footnote varchar, min_year_acquired varchar);

// --COPY food_nutrient FROM 'D:\.FOR PRODUCTIVITY\Manoc\FoodData Central\excel\food_nutrient.csv' with csv header;

// --

// --create table food(fdc_id int, data_type varchar, description varchar, food_category_id varchar, publication_date varchar);

// --COPY food FROM 'D:\.FOR PRODUCTIVITY\Manoc\FoodData Central\excel\food.csv' with csv header;

// --
package setup

var backup = `

CREATE TABLE ingredient(
    id serial primary key,
    name varchar not null UNIQUE,
    name_ph varchar DEFAULT '',
    name_brand varchar DEFAULT '',
    date_created date,
    barcode varchar unique,
    thumbnail_image_link varchar,
    ingredient_desc varchar default '',
    category_id int,
    FOREIGN KEY(category_id) REFERENCES edible_category(id)
);
CREATE TABLE ingredient_variant(
    id serial primary key,
    name varchar not null,
    ingredient_id int NOT NULL
);
CREATE TABLE ingredient_cook_type(
    id serial primary key,
    name varchar not null,
    ingredient_id int NOT NULL,
    ingredient_variant_id int NOT NULL,
    nutrient_id int not null UNIQUE,
    FOREIGN KEY(ingredient_id) REFERENCES ingredient(id),
    FOREIGN KEY(ingredient_variant_id) REFERENCES ingredient_variant(id),
    FOREIGN KEY(nutrient_id) REFERENCES nutrient(id) ON DELETE cascade
);
create table ingredient_image(
    id serial primary key,
    ingredient_id int not null,
    name_file varchar not NULL,
    amount float4 not NULL,
    amount_unit varchar(4) not NULL,
    amount_unit_desc varchar(40) not NULL
);

create table nutrient(
    id serial primary key,
    amount float4 not NULL,
    amount_unit varchar(4) not NULL,
    amount_unit_desc varchar(40) not NULL,
    serving_size float4 default 0,
    calories float4 not NULL,
    protein float4 not NULL,
    carbs float4 not NULL,
    fats float4 not null,
    trans_fat float4,
    saturated_fat float4,
    sugars float4,
    fiber float4,
    sodium float4,
    iron float4,
    calcium float4,
);

create table edible_category(
    id int primary key,
    name varchar not null
);

-- For Text Search

create text search dictionary simple_tsd (
    template = pg_catalog.simple,
    stopwords = english
);

create text search configuration simple_tsc (copy = simple);

alter text search configuration simple_tsc alter mapping for asciiword
with simple_tsd;

alter table food add column search_food tsvector generated always as (
    setweight(to_tsvector('english',coalesce(food.name, '')),'A') ||
    setweight(to_tsvector('english',coalesce(food.name_ph, '')),'B') ||
	setweight(to_tsvector('english',coalesce(food.name_brand, '')),'C')
) STORED;

create index search_food_idx on food using GIN(search_food);

--sample query (full text)

select
    food.name,
    food.name_brand,
    food_nutrient.calories,
    food_nutrient.carbs,
    food_nutrient.fats,
    food_nutrient.protein,
    ts_rank_cd(
        search_food,
        to_tsquery('english', 'chicken')
    ) as ranking
from food
    JOIN food_nutrient ON food.food_nutrient_id = food_nutrient.id
where
    search_food @@to_tsquery('english', 'chicken')
    and name_brand = 'USDA'
    and food_category_id = 6 
order by ranking desc;

select
    food.name,
    food.name_brand,
    food_nutrient.calories,
    food_nutrient.carbs,
    food_nutrient.fats,
    food_nutrient.protein,
    ts_rank_cd(
        search_food,
        to_tsquery('english', 'chicken')
    ) as ranking
from food
    JOIN food_nutrient ON food.food_nutrient_id = food_nutrient.id
where
    search_food @@to_tsquery(
        'english',
        'chick:* & broiler:* & fryer:*'
    )
    and name_brand = 'USDA'
    and food_category_id = 6 
    and name not like '%rotisserie%'
order by ranking desc;

-- Deprecated

--alter table food add column search_food tsvector;

--update "food" set search_food = to_tsvector(food.name || ' ' || coalesce(food.name_ph, '') || ' ' || coalesce(food.name_brand, ''));

---- Function for creating the ts vector

--CREATE FUNCTION search_food_vector_func() RETURNS TRIGGER AS $$

--BEGIN

--    NEW.search_food := to_tsvector(name || ' ' || coalesce(name_ph, '') || ' ' || coalesce(name_brand, ''));

--    RETURN NEW;

--END;

--$$ LANGUAGE plpgsql;

--CREATE TRIGGER search_food_vector_update

--BEFORE INSERT OR UPDATE ON food

--FOR EACH ROW

--EXECUTE FUNCTION search_food_vector_func();

-- select * from food where search_food @@ to_tsquery('chicken') ORDER By score DESC;

create table
    food(
        id serial primary key,
        name varchar not null UNIQUE,
        name_ph varchar DEFAULT '',
        name_brand varchar DEFAULT '',
        date_created date,
        barcode varchar unique,
        thumbnail_image_link varchar,
        food_desc varchar default '',
        food_nutrient_id int not null UNIQUE,
        food_brand_type_id int not null,
        food_category_id int,
        food_brand_id uuid not null,
        removed bool DEFAULT FALSE NOT NULL,
        FOREIGN KEY(food_brand_type_id) REFERENCES food_brand_type(id),
        FOREIGN KEY(food_nutrient_id) REFERENCES food_nutrient(id) ON DELETE cascade,
        FOREIGN KEY(food_category_id) REFERENCES food_category(id),
        FOREIGN KEY(food_brand_id) REFERENCES food_brand(id)
    );
create table
    food_nutrient(
        id serial primary key,
        amount float4 not NULL,
        amount_unit varchar(4) not NULL,
        amount_unit_desc varchar(40) not NULL,
        serving_size float4 default 0,
        calories float4 not NULL,
        protein float4 not NULL,
        carbs float4 not NULL,
        fats float4 not null,
        trans_fat float4,
        saturated_fat float4,
        sugars float4,
        fiber float4,
        sodium float4,
        iron float4,
        calcium float4
    );

create table
    food_image(
        id serial primary key,
        food_id int not null,
        name_file varchar not NULL,
        amount float4 not NULL,
        amount_unit varchar(4) not NULL,
        amount_unit_desc varchar(40) not NULL
    );

create table
    food_category(
        id int primary key,
        name varchar not null
    );

create table
    food_rewards(
        id serial primary key,
        food_id int not null,
        coins int not null,
        xp int not null
    );

create table
    food_brand(
        id uuid primary key,
        name varchar not null,
        brand_desc varchar,
        thumbnail_image_link varchar,
        cover_image_link varchar,
        profile_image_link varchar,
        food_brand_type_id int not null,
        FOREIGN KEY(food_brand_type_id) REFERENCES food_brand_type(id)
    );

create table
    food_brand_type(
        id serial primary key,
        name varchar not null,
        brand_type_desc varchar
    );

insert into
    food_brand_type(name)
values ('generic'), ('commercial');

insert into
    food_brand(id, name, food_brand_type_id)
values (uuid_generate_v4(), 'iFNRI', 1), (uuid_generate_v4(), 'USDA', 1);

insert into edible_category(id, name)
    values 
        (1, 'cereals'),
        (2, 'starchy roots, and tubers'),
        (3, 'nuts, dried beans, and seeds'),
        (4, 'vegetables'),
        (5, 'fruits'),
        (6, 'meats'),
        (7, 'seafood'),
        (8, 'eggs'),
        (9, 'dairy'),
        (10, 'fats and oils'),
        (11, 'sweets'),
        (12, 'spices and herbs'),
        (13, 'alcoholic beverages'),
        (14, 'non-alcoholic beverages'),
        (15, 'baby foods'),
        (16, 'soups, sauces, and gravies'),
        (17, 'miscellaneous'),
        (18, 'branded');

create table
    recipe(
        id serial primary key,
        owner_id uuid not NULL,
        recipe_nutrient_id int not NULL,
        name varchar not NULL,
        main_image_link varchar,
        thumbnail_image_link varchar,
        prep_mins int not NULL,
        servings int not NULL,
        likes int,
        date_created timestamp,
        date_updated timestamp,
        FOREIGN KEY(owner_id) REFERENCES account(id),
        FOREIGN KEY(recipe_nutrient_id) REFERENCES recipe_nutrient(id)
    );

create table
    recipe_nutrient(
        id serial primary key,
        amount float4 not NULL,
        amount_unit varchar(4) not NULL,
        amount_unit_desc varchar(40) not NULL,
        serving_size float4 default 0,
        calories float4 not NULL,
        protein float4 not NULL,
        carbs float4 not NULL,
        fats float4 not null,
        trans_fat float4 default 0,
        saturated_fat float4 default 0,
        sugars float4 default 0,
        sodium float4 default 0
    );

create table
    recipe_ingredient(
        id serial primary key,
        food_id int not NULL,
        amount float4 not NULL,
        amount_unit varchar(4) not NULL,
        amount_unit_desc varchar(40) not NULL,
        serving_size float4,
        FOREIGN KEY(food_id) REFERENCES food(id)
    );

create table
    recipe_step(
        id serial primary key,
        recipe_id int not NULL,
        step_order int not NULL,
        step_desc varchar not NULL
    );

create table account(
    id uuid primary key,
    username varchar unique NOT NULL,
    password varchar NOT NULL,
    name_first varchar,
    name_last varchar,
    phone_number varchar,
    date_updated timestamp,
    date_created timestamp,
    account_type_id uuid NOT NULL,
    account_vitals_id uuid NOT NULL,
    account_profile_id uuid NOT NULL,
    measure_unit_id uuid NOT NULL,
    FOREIGN KEY(account_type_id) REFERENCES account_type(id),
    FOREIGN KEY(account_vitals_id) REFERENCES account_vitals(id),
    FOREIGN KEY(account_profile_id) REFERENCES account_profile(id),
    FOREIGN KEY(measure_unit_id) REFERENCES measure_unit(id)
);

create table account_type(
    id uuid primary key,
    name varchar unique NOT NULL,
    account_type_desc varchar
);

create table account_vitals(
    id uuid primary key,
    account_id uuid not NULL unique,
    weight int2 NOT NULL,
    height int2 NOT NULL,
    birthday date NOT NULL,
    sex varchar NOT null,
    activity_lvl_id uuid NOT NULL,
    diet_plan_id uuid NOT NULL,
    --    FOREIGN KEY(account_id) REFERENCES account(id),
    FOREIGN KEY(activity_lvl_id) REFERENCES activity_lvl(id),
    FOREIGN KEY(diet_plan_id) REFERENCES diet_plan(id)
);

create table account_weight_changes(
    id serial primary key,
    account_id uuid not NULL unique,
    weight int2 NOT NULL,
    date_created date
);

create table account_profile(
    id uuid primary key,
    account_id uuid not null UNIQUE,
    account_image_link varchar,
    account_title varchar --    FOREIGN KEY(account_id) REFERENCES account(id),
);

create table account_items(
    id serial primary key,
    account_id uuid not NULL,
    game_item_id int not null,
    FOREIGN KEY(game_item_id) REFERENCES game_item(id)
);

create table account_game_stat(
    id serial primary key,
    account_id uuid not null UNIQUE,
    coins int,
    xp int
);

insert into account_type(id, name, account_type_desc)
    values
    (uuid_generate_v4(), 'user', ''),
    (uuid_generate_v4(), 'dietician', ''),
    (uuid_generate_v4(), 'brand_owner', '');

create table game_item(
    id serial primary key,
    name varchar,
    price int,
    commentary varchar,
    game_item_desc varchar,
    game_item_image_link varchar
);

create table measure_unit(
    id uuid primary key,
    name varchar,
    wt_solid_unit_small varchar(4) NOT NULL,
    wt_solid_desc_small varchar(40) NOT NULL,
    wt_solid_unit_medium varchar(4) NOT NULL,
    wt_solid_desc_medium varchar(40) NOT NULL,
    wt_solid_unit_large varchar(4) NOT NULL,
    wt_solid_desc_large varchar(40) NOT NULL,
    wt_liquid_unit_small varchar(4) NOT NULL,
    wt_liquid_desc_small varchar(40) NOT NULL,
    wt_liquid_unit_medium varchar(4) NOT NULL,
    wt_liquid_desc_medium varchar(40) NOT null,
    wt_liquid_unit_large varchar(4) NOT NULL,
    wt_liquid_desc_large varchar(40) NOT null
);

create table diet_plan(
    id uuid primary key,
    name varchar UNIQUE,
    main_image_link varchar,
    background_color varchar,
    diet_plan_desc varchar,
    calorie_percentage int,
    protein_percentage int,
    fats_percentage int,
    carbs_percentage int
);

create table activity_lvl(
    id uuid primary key,
    name varchar,
    main_image_link varchar,
    background_color varchar,
    activity_lvl_desc varchar,
    bmr_multiplier float
);

create table daily_nutrients(
    id serial primary key,
    account_id uuid not NULL,
    date_created date not NULL,
    calories float4 not NULL,
    protein float4 not NULL,
    carbs float4 not NULL,
    fats float4 not NULL,
    max_calories float4 not NULL,
    max_protein float4 not NULL,
    max_carbs float4 not NULL,
    max_fats float4 not NULL,
    activity_lvl_id uuid NOT NULL,
    diet_plan_id uuid NOT NULL,
    FOREIGN KEY(activity_lvl_id) REFERENCES activity_lvl(id),
    FOREIGN KEY(diet_plan_id) REFERENCES diet_plan(id)
);

create table intake(
    id serial primary key,
    account_id uuid not NULL,
    date_created timestamp,
    amount float4 not NULL,
    amount_unit varchar(4) not NULL,
    amount_unit_desc varchar(40) not NULL,
    serving_size float4 default 0,
    food_id int,
    FOREIGN KEY(food_id) REFERENCES food(id)
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

insert into
    measure_unit(
        id,
        name,
        wt_solid_unit_small,
        wt_solid_desc_small,
        wt_solid_unit_medium,
        wt_solid_desc_medium,
        wt_solid_unit_large,
        wt_solid_desc_large,
        wt_liquid_unit_small,
        wt_liquid_desc_small,
        wt_liquid_unit_medium,
        wt_liquid_desc_medium,
        wt_liquid_unit_large,
        wt_liquid_desc_large
    )
values (
        uuid_generate_v4(),
        'metric',
        'mg',
        'milligram',
        'g',
        'gram',
        'kg',
        'kilogram',
        '',
        '',
        'ml',
        'millilitre',
        'l',
        'litre'
    );

insert into diet_plan(id, name, calorie_percentage, protein_percentage, fats_percentage, carbs_percentage)
    values
    (uuid_generate_v4(), 'mild weight loss', 90, 10, 30, 60),
    (uuid_generate_v4(), 'mild fat loss (low carb)', 90, 20, 30, 50),
    (uuid_generate_v4(), 'mild fat loss (low fat)', 90, 20, 20, 60),
    (uuid_generate_v4(), 'weight loss', 80, 12, 30, 58),
    (uuid_generate_v4(), 'fat loss (low carb)', 80, 22, 30, 48),
    (uuid_generate_v4(), 'fat loss (low fat)', 80, 22, 18, 60),
    (uuid_generate_v4(), 'more moderate weight loss', 70, 15, 30, 55),
    (uuid_generate_v4(), 'more moderate fat loss (low carb)', 70, 25, 30, 45),
    (uuid_generate_v4(), 'more moderate fat loss (low fat)', 70, 25, 15, 60),
    (uuid_generate_v4(), 'extreme weight loss', 60, 20, 25, 55),
    (uuid_generate_v4(), 'extreme fat loss (low carb)', 60, 30, 30, 40),
    (uuid_generate_v4(), 'extreme fat loss (low fat)', 60, 30, 15, 55);

insert into activity_lvl(id, name, bmr_multiplier)
    values
    (uuid_generate_v4(), 'inactive', 1.2),
    (uuid_generate_v4(), 'lightly active', 1.375),
    (uuid_generate_v4(), 'moderately active', 1.465),
    (uuid_generate_v4(), 'active', 1.55),
    (uuid_generate_v4(), 'very active', 1.725),
    (uuid_generate_v4(), 'extremely active', 1.9);

SELECT
    account_vitals.account_id,
    account_vitals.weight,
    account_vitals.height,
    account_vitals.birthday,
    account_vitals.sex,
    account_vitals.activity_lvl_id,
    account_vitals.diet_plan_id,
    activity_lvl.name,
    activity_lvl.bmr_multiplier,
    diet_plan.name,
    diet_plan.calorie_percentage,
    diet_plan.protein_percentage,
    diet_plan.fats_percentage,
    diet_plan.carbs_percentage
FROM account_vitals
    JOIN activity_lvl ON account_vitals.activity_lvl_id = activity_lvl.id
    JOIN diet_plan ON account_vitals.diet_plan_id = diet_plan.id
WHERE
    account_vitals.account_id = '898f8e6c-817e-4605-af14-5b437c58bc86';

insert into
    account_game_stat (account_id, coins, xp)
values (
        '898f8e6c-817e-4605-af14-5b437c58bc86',
        0,
        0
    );

ALTER TABLE food MODIFY COLUMN name_brand VARCHAR() DEFAULT '';

update food_nutrient set serving_size = 0 where serving_size is null;

SELECT name FROM food;`
