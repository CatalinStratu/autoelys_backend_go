-- Mock Vehicles Data for Testing
-- This script creates sample vehicle listings for the existing user

-- First, let's ensure we have the user ID (assuming user ID 3 exists)
SET @user_id = 3;

-- Insert mock vehicles
INSERT INTO vehicles (
    user_id, uuid, slug, title, category, description, price, currency, negotiable,
    person_type_id, brand, model, engine_capacity, power_hp,
    fuel_type_id, body_type_id, kilometers, color, year, number_of_keys,
    condition_id, transmission_id, steering_id, registered,
    city, contact_name, email, phone
) VALUES
-- Vehicle 1: BMW 320d
(
    @user_id,
    '550e8400-e29b-41d4-a716-446655440000',
    'bmw-320d-2019-xdrive-impecabil',
    'BMW 320d 2019 xDrive - Impecabil',
    'autoturisme',
    'BMW Seria 3 320d in stare excelenta. Motor diesel 2.0L, tractiune integrala xDrive. Intretinere completa la reprezentanta. Masina este ca noua, fara accident.',
    18500.00,
    'lei',
    true,
    1, -- persoana_fizica
    'BMW',
    '320d',
    1995,
    190,
    2, -- motorina
    1, -- sedan
    85000,
    'Negru',
    2019,
    2,
    1, -- utilizat
    2, -- automata
    1, -- stanga
    true,
    'București',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 2: Audi A4
(
    @user_id,
    '660e8400-e29b-41d4-a716-446655440001',
    'audi-a4-2020-quattro-s-line',
    'Audi A4 2020 Quattro S-Line',
    'autoturisme',
    'Audi A4 in configuratie S-Line, dotari premium. Quattro tractiune integrala, interior piele, navigatie, camera 360. Service la zi, un singur proprietar.',
    22000.00,
    'lei',
    false,
    1, -- persoana_fizica
    'Audi',
    'A4',
    1984,
    204,
    2, -- motorina
    1, -- sedan
    62000,
    'Gri',
    2020,
    2,
    1, -- utilizat
    2, -- automata
    1, -- stanga
    true,
    'Cluj-Napoca',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 3: Mercedes-Benz C-Class
(
    @user_id,
    '770e8400-e29b-41d4-a716-446655440002',
    'mercedes-benz-c200-2021-amg',
    'Mercedes-Benz C200 2021 AMG',
    'autoturisme',
    'Mercedes C-Class C200 AMG Line, dotari complete: scaune incalzite/ventilate, panoramic, LED Matrix, distronic, keyless, head-up display. Garantie Mercedes.',
    28000.00,
    'lei',
    true,
    2, -- firma
    'Mercedes-Benz',
    'C200',
    1991,
    204,
    1, -- benzina
    1, -- sedan
    45000,
    'Alb',
    2021,
    2,
    1, -- utilizat
    2, -- automata
    1, -- stanga
    true,
    'Timișoara',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 4: Volkswagen Golf
(
    @user_id,
    '880e8400-e29b-41d4-a716-446655440003',
    'vw-golf-8-2022-gti',
    'VW Golf 8 2022 GTI',
    'autoturisme',
    'Volkswagen Golf 8 GTI, motor 2.0 TSI 245 CP. Masina sport, performanta exceptionala. Pachet Digital Cockpit Pro, DCC, scaune sport Alcantara.',
    26500.00,
    'lei',
    false,
    1, -- persoana_fizica
    'Volkswagen',
    'Golf GTI',
    1984,
    245,
    1, -- benzina
    8, -- hatchback
    28000,
    'Roșu',
    2022,
    2,
    1, -- utilizat
    2, -- automata
    1, -- stanga
    true,
    'Brașov',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 5: Toyota RAV4
(
    @user_id,
    '990e8400-e29b-41d4-a716-446655440004',
    'toyota-rav4-2023-hybrid-awd',
    'Toyota RAV4 2023 Hybrid AWD',
    'autoturisme',
    'Toyota RAV4 Hybrid, tractiune AWD, consum mic, fiabilitate maxima. Dotari: JBL Premium Sound, panoramic, scaune piele, senzori 360, adaptive cruise control.',
    32000.00,
    'lei',
    true,
    1, -- persoana_fizica
    'Toyota',
    'RAV4',
    2487,
    222,
    4, -- hibrid
    2, -- suv
    15000,
    'Argintiu',
    2023,
    2,
    1, -- utilizat
    2, -- automata
    1, -- stanga
    true,
    'Constanța',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 6: Ford Mustang
(
    @user_id,
    'aa0e8400-e29b-41d4-a716-446655440005',
    'ford-mustang-2020-gt-v8',
    'Ford Mustang 2020 GT V8',
    'autoturisme',
    'Ford Mustang GT 5.0 V8 450 CP! Masina muscle car autentica americana. Sunet V8 de exceptie, acceleratii brutale. Scaune Recaro, sistem audio premium.',
    35000.00,
    'lei',
    false,
    1, -- persoana_fizica
    'Ford',
    'Mustang GT',
    5038,
    450,
    1, -- benzina
    4, -- coupe
    42000,
    'Portocaliu',
    2020,
    2,
    1, -- utilizat
    2, -- automata
    1, -- stanga
    true,
    'București',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 7: Dacia Duster
(
    @user_id,
    'bb0e8400-e29b-41d4-a716-446655440006',
    'dacia-duster-2021-4x4-prestige',
    'Dacia Duster 2021 4x4 Prestige',
    'autoturisme',
    'Dacia Duster 4x4 in versiune Prestige, dotari complete pentru clasa sa. Tractiune 4x4, navigatie, camera marsarier, senzori parcare. Ideal offroad usor.',
    14500.00,
    'lei',
    true,
    1, -- persoana_fizica
    'Dacia',
    'Duster',
    1461,
    115,
    2, -- motorina
    2, -- suv
    68000,
    'Portocaliu',
    2021,
    2,
    1, -- utilizat
    1, -- manuala
    1, -- stanga
    true,
    'Iași',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 8: Porsche 911
(
    @user_id,
    'cc0e8400-e29b-41d4-a716-446655440007',
    'porsche-911-carrera-2022',
    'Porsche 911 Carrera 2022',
    'autoturisme',
    'Porsche 911 Carrera, iconic sports car. Motor boxer 3.0L 385 CP, 0-100 km/h in 4.2 sec. PASM, Sport Chrono, scaune sport adaptive, carbon interior.',
    95000.00,
    'lei',
    false,
    2, -- firma
    'Porsche',
    '911 Carrera',
    2981,
    385,
    1, -- benzina
    4, -- coupe
    18000,
    'Albastru',
    2022,
    2,
    1, -- utilizat
    2, -- automata
    1, -- stanga
    true,
    'București',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 9: Renault Clio
(
    @user_id,
    'dd0e8400-e29b-41d4-a716-446655440008',
    'renault-clio-2023-nou-intens',
    'Renault Clio 2023 NOU - Intens',
    'autoturisme',
    'Renault Clio 5 NOU, versiune Intens. Masina mica, economica, perfecta pentru oras. Dotari: ecran tactil 9.3", climate control, camera, LED lights.',
    13500.00,
    'lei',
    false,
    2, -- firma
    'Renault',
    'Clio',
    999,
    90,
    1, -- benzina
    8, -- hatchback
    0,
    'Alb',
    2023,
    2,
    2, -- nou
    1, -- manuala
    1, -- stanga
    true,
    'Pitești',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
),

-- Vehicle 10: Tesla Model 3
(
    @user_id,
    'ee0e8400-e29b-41d4-a716-446655440009',
    'tesla-model-3-2023-long-range',
    'Tesla Model 3 2023 Long Range',
    'autoturisme',
    'Tesla Model 3 Long Range, autonomie 614 km, autopilot, 0-100 in 4.4 sec. Interior premium, ecran 15", sticla panoramica, sistem audio premium. Full electric.',
    42000.00,
    'lei',
    true,
    1, -- persoana_fizica
    'Tesla',
    'Model 3',
    0,
    283,
    3, -- electric
    1, -- sedan
    32000,
    'Negru',
    2023,
    0,
    1, -- utilizat
    2, -- automata
    1, -- stanga
    true,
    'Cluj-Napoca',
    'Ion Popescu',
    'john@example.com',
    '0721123456'
);

-- Add some sample images for the first 3 vehicles
INSERT INTO vehicle_images (vehicle_id, image_url)
SELECT id, CONCAT('/uploads/vehicles/sample-', id, '-1.jpg') FROM vehicles WHERE uuid = '550e8400-e29b-41d4-a716-446655440000';

INSERT INTO vehicle_images (vehicle_id, image_url)
SELECT id, CONCAT('/uploads/vehicles/sample-', id, '-2.jpg') FROM vehicles WHERE uuid = '550e8400-e29b-41d4-a716-446655440000';

INSERT INTO vehicle_images (vehicle_id, image_url)
SELECT id, CONCAT('/uploads/vehicles/sample-', id, '-1.jpg') FROM vehicles WHERE uuid = '660e8400-e29b-41d4-a716-446655440001';

INSERT INTO vehicle_images (vehicle_id, image_url)
SELECT id, CONCAT('/uploads/vehicles/sample-', id, '-1.jpg') FROM vehicles WHERE uuid = '770e8400-e29b-41d4-a716-446655440002';
