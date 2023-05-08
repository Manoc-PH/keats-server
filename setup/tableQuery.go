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
    name varchar not null UNIQUE, 
    name_ph varchar UNIQUE
);
CREATE TABLE ingredient_subvariant(
    id serial primary key,
    name varchar not null UNIQUE, 
    name_ph varchar UNIQUE
);
CREATE TABLE ingredient_nutrient(
    id serial primary key,
    ingredient_id int NOT NULL,
    ingredient_variant_id int NOT NULL,
    ingredient_subvariant_id int NOT NULL,
    nutrient_id int not null UNIQUE,
    FOREIGN KEY(ingredient_id) REFERENCES ingredient(id),
    FOREIGN KEY(ingredient_variant_id) REFERENCES ingredient_variant(id),
    FOREIGN KEY(ingredient_subvariant_id) REFERENCES ingredient_subvariant(id),
    FOREIGN KEY(nutrient_id) REFERENCES nutrient(id) ON DELETE cascade
);
CREATE TABLE ingredient_image(
    id serial primary key,
    ingredient_mapping_id int not null,
    name_file varchar not NULL,
    amount float4 not NULL,
    amount_unit varchar(4) not NULL,
    amount_unit_desc varchar(40) not NULL
);

CREATE TABLE nutrient(
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

CREATE TABLE edible_category(
    id int primary key,
    name varchar not null
);

-- VIEW
CREATE VIEW ingredient_details AS
    SELECT 
        ingredient.id,
        ingredient.name,
        ingredient.name_ph,
        ingredient.name_brand,
        ingredient_variant.name as variant_name,
        ingredient_subvariant.name as subvariant_name,
        nutrient.amount,
        nutrient.amount_unit,
        nutrient.calories,
        nutrient.protein,
        nutrient.carbs,
        nutrient.fats
    FROM ingredient_mapping 
    JOIN ingredient on ingredient_mapping.ingredient_id = ingredient.id
    JOIN ingredient_variant on ingredient_mapping.ingredient_variant_id = ingredient_variant.id 
    JOIN ingredient_subvariant on ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id 
    JOIN nutrient on ingredient_mapping.nutrient_id = nutrient.id 



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

select * from food 
where 
    search_food @@to_tsquery(
        'english',
        'Alcoholic & beverage'
    ) 
	and removed = false
    and name_brand = 'USDA'
    and food_category_id = 14
order by name desc;

select name, removed from food 
where removed = false
    and name_brand = 'USDA'
    and food_category_id = 17
order by name asc;

SELECT 
    ingredient.name AS ingredient_name,
    ingredient_variant.name as part,
    ingredient_subvariant.name as cook_type,
    nutrient.calories,
    nutrient.carbs,
    nutrient.fats,
    nutrient.protein
FROM ingredient_mapping
    JOIN ingredient on ingredient_mapping.ingredient_id = ingredient.id
    JOIN ingredient_variant on ingredient_mapping.ingredient_variant_id = ingredient_variant.id
    JOIN ingredient_subvariant on ingredient_mapping.ingredient_subvariant_id = ingredient_subvariant.id
    JOIN nutrient on ingredient_mapping.nutrient_id = nutrient.id;
where ingredient.category_id = 4;
 
select * from FOOD where removed = false and name_brand = 'USDA' ;
--delete from ingredient_variant where ingredient_id = 6;
--update food set removed = false where food_category_id = 6 ; 

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

-- create table
--     food(
--         id serial primary key,
--         name varchar not null UNIQUE,
--         name_ph varchar DEFAULT '',
--         name_brand varchar DEFAULT '',
--         date_created date,
--         barcode varchar unique,
--         thumbnail_image_link varchar,
--         food_desc varchar default '',
--         food_nutrient_id int not null UNIQUE,
--         food_brand_type_id int not null,
--         food_category_id int,
--         food_brand_id uuid not null,
--         removed bool DEFAULT FALSE NOT NULL,
--         FOREIGN KEY(food_brand_type_id) REFERENCES food_brand_type(id),
--         FOREIGN KEY(food_nutrient_id) REFERENCES food_nutrient(id) ON DELETE cascade,
--         FOREIGN KEY(food_category_id) REFERENCES food_category(id),
--         FOREIGN KEY(food_brand_id) REFERENCES food_brand(id)
--     );
-- create table
--     food_nutrient(
--         id serial primary key,
--         amount float4 not NULL,
--         amount_unit varchar(4) not NULL,
--         amount_unit_desc varchar(40) not NULL,
--         serving_size float4 default 0,
--         calories float4 not NULL,
--         protein float4 not NULL,
--         carbs float4 not NULL,
--         fats float4 not null,
--         trans_fat float4,
--         saturated_fat float4,
--         sugars float4,
--         fiber float4,
--         sodium float4,
--         iron float4,
--         calcium float4
--     );

-- create table
--     food_rewards(
--         id serial primary key,
--         food_id int not null,
--         coins int not null,
--         xp int not null
--     );

CREATE TABLE food (
    id serial primary key,
    name varchar not null UNIQUE,
    name_ph varchar DEFAULT '',
    name_brand varchar DEFAULT '',
    date_created date,
    barcode varchar unique,
    thumbnail_image_link varchar,
    food_desc varchar default '', 
    brand_type_id int not null,
    category_id int not null,
    brand_id uuid not null,
    FOREIGN KEY(brand_type_id) REFERENCES edible_brand_type(id),
    FOREIGN KEY(category_id) REFERENCES edible_category(id),
    FOREIGN KEY(brand_id) REFERENCES edible_brand(id)
);
CREATE TABLE food_ingredient (
    id serial primary key,
    ingredient_mapping_id serial,
    amount float4 not NULL,
    amount_unit varchar(4) not NULL,
    amount_unit_desc varchar(40) not NULL,
    serving_size float4 default 0,
    FOREIGN KEY(ingredient_mapping_id) REFERENCES ingredient_mapping(id)
);
CREATE TABLE food_image(
    id serial primary key,
    food_id int not null,
    name_file varchar not NULL,
    amount float4 not NULL,
    amount_unit varchar(4) not NULL,
    amount_unit_desc varchar(40) not NULL
);
CREATE TABLE edible_category(
    id int primary key,
    name varchar not null
);

CREATE TABLE edible_brand(
    id uuid primary key,
    name varchar not null,
    brand_desc varchar,
    thumbnail_image_link varchar,
    cover_image_link varchar,
    profile_image_link varchar,
    brand_type_id int not null,
    FOREIGN KEY(brand_type_id) REFERENCES edible_brand_type(id)
);
CREATE TABLE edible_brand_type(
    id serial primary key,
    name varchar not null,
    brand_type_desc varchar
);

CREATE TABLE food_intake(
    id serial primary key,
    food_id int NOT NULL,
    FOREIGN KEY(food_id) REFERENCES food(id)
);
CREATE TABLE food_intake_mapping(
    id serial primary key,
    food_intake_id int NOT NULL,
    ingredient_mapping_id int NOT NULL,
    amount float4 not NULL,
    amount_unit varchar(4) not NULL,
    FOREIGN KEY(food_intake_id) REFERENCES food_intake(id),
    FOREIGN KEY(ingredient_mapping_id) REFERENCES ingredient_mapping(id)
);
 
CREATE TABLE intake(
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

insert into edible_brand_type(name)
values ('generic'), ('commercial');

insert into edible_brand(id, name, edible_brand_type_id)
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

create table account(
    id uuid primary key,
    username varchar unique NOT NULL,
    password varchar NOT NULL,
    name_first varchar,
    name_last varchar,
    phone_number varchar,
    date_updated timestamp,
    date_created timestamp,
    account_vitals_id uuid NOT NULL,
    account_profile_id uuid NOT NULL,
    measure_unit_id uuid NOT NULL,
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
    account_title varchar,
    account_type_id uuid NOT NULL,
    FOREIGN KEY(account_type_id) REFERENCES account_type(id)
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

INSERT INTO account_type (id, name)
	values
		('4c3c69b0-2eae-4b3c-80e1-619f4718d272', 'consumer'),
		('7d3f6af5-acd7-49e4-b968-692b7301fa6c', 'admin'),
		('a65ddf3e-9d55-4da9-b695-69d0aaeeedab', 'business');

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

insert into activity_lvl(id, name, bmr_multiplier, activity_lvl_desc)
    values
    (uuid_generate_v4(), 'inactive', 1.2, 'This is characterized by low levels of physical activity and a lot of sitting or lying down throughout the day.'),
    (uuid_generate_v4(), 'lightly active', 1.375, 'Increased heart rate and breathing for at least about 15-30 minutes, 1-3 days a week.'),
    (uuid_generate_v4(), 'moderately active', 1.465, 'Increased heart rate and breathing for at least about 15-30 minutes, 4-5 days a week.'),
    (uuid_generate_v4(), 'active', 1.55, 'Increased heart rate and breathing for at least about 45-120 minutes, 3-4 days a week.'),
    (uuid_generate_v4(), 'very active', 1.725, 'Increased heart rate and breathing for at least about 2 hours or more, 6-7 days a week.'),
    (uuid_generate_v4(), 'extremely active', 1.9, 'Increased heart rate and breathing for at least about 2 hours or more, everyday.');

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

insert into account_game_stat (account_id, coins, xp)
	values ( '898f8e6c-817e-4605-af14-5b437c58bc86', 0, 0 );

ALTER TABLE food MODIFY COLUMN name_brand VARCHAR() DEFAULT '';

update food_nutrient set serving_size = 0 where serving_size is null;

SELECT name FROM food;`
