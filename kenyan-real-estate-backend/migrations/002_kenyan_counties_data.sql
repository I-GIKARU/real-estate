-- Migration: 002_kenyan_counties_data.sql
-- Populate counties and sub-counties with Kenyan data

-- Insert Kenyan counties
INSERT INTO counties (name, code) VALUES
('Baringo', 'BRG'),
('Bomet', 'BMT'),
('Bungoma', 'BGM'),
('Busia', 'BSA'),
('Elgeyo-Marakwet', 'EGM'),
('Embu', 'EMB'),
('Garissa', 'GRS'),
('Homa Bay', 'HMB'),
('Isiolo', 'ISL'),
('Kajiado', 'KJD'),
('Kakamega', 'KKG'),
('Kericho', 'KRC'),
('Kiambu', 'KMB'),
('Kilifi', 'KLF'),
('Kirinyaga', 'KRG'),
('Kisii', 'KSI'),
('Kisumu', 'KSM'),
('Kitui', 'KTI'),
('Kwale', 'KWL'),
('Laikipia', 'LKP'),
('Lamu', 'LAM'),
('Machakos', 'MCK'),
('Makueni', 'MKN'),
('Mandera', 'MND'),
('Marsabit', 'MSB'),
('Meru', 'MRU'),
('Migori', 'MGR'),
('Mombasa', 'MSA'),
('Murang''a', 'MRG'),
('Nairobi', 'NRB'),
('Nakuru', 'NKR'),
('Nandi', 'NND'),
('Narok', 'NRK'),
('Nyamira', 'NYM'),
('Nyandarua', 'NND'),
('Nyeri', 'NYR'),
('Samburu', 'SMB'),
('Siaya', 'SYA'),
('Taita-Taveta', 'TTT'),
('Tana River', 'TNR'),
('Tharaka-Nithi', 'THN'),
('Trans Nzoia', 'TNZ'),
('Turkana', 'TRK'),
('Uasin Gishu', 'UGS'),
('Vihiga', 'VHG'),
('Wajir', 'WJR'),
('West Pokot', 'WPK');

-- Insert major sub-counties for key counties (focusing on major urban areas)

-- Nairobi County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'NRB'), 'Westlands'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Dagoretti North'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Dagoretti South'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Langata'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Kibra'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Roysambu'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Kasarani'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Ruaraka'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Embakasi South'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Embakasi North'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Embakasi Central'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Embakasi East'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Embakasi West'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Makadara'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Kamukunji'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Starehe'),
((SELECT id FROM counties WHERE code = 'NRB'), 'Mathare');

-- Mombasa County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'MSA'), 'Changamwe'),
((SELECT id FROM counties WHERE code = 'MSA'), 'Jomba'),
((SELECT id FROM counties WHERE code = 'MSA'), 'Kisauni'),
((SELECT id FROM counties WHERE code = 'MSA'), 'Nyali'),
((SELECT id FROM counties WHERE code = 'MSA'), 'Likoni'),
((SELECT id FROM counties WHERE code = 'MSA'), 'Mvita');

-- Kiambu County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'KMB'), 'Thika Town'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Ruiru'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Juja'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Kiambu'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Kiambaa'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Kabete'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Kikuyu'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Limuru'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Githunguri'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Lari'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Gatundu South'),
((SELECT id FROM counties WHERE code = 'KMB'), 'Gatundu North');

-- Nakuru County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'NKR'), 'Nakuru Town East'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Nakuru Town West'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Bahati'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Subukia'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Rongai'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Kuresoi South'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Kuresoi North'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Gilgil'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Naivasha'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Molo'),
((SELECT id FROM counties WHERE code = 'NKR'), 'Njoro');

-- Uasin Gishu County Sub-Counties (Eldoret)
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'UGS'), 'Ainabkoi'),
((SELECT id FROM counties WHERE code = 'UGS'), 'Kapseret'),
((SELECT id FROM counties WHERE code = 'UGS'), 'Kesses'),
((SELECT id FROM counties WHERE code = 'UGS'), 'Moiben'),
((SELECT id FROM counties WHERE code = 'UGS'), 'Soy'),
((SELECT id FROM counties WHERE code = 'UGS'), 'Turbo');

-- Kisumu County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'KSM'), 'Kisumu Central'),
((SELECT id FROM counties WHERE code = 'KSM'), 'Kisumu East'),
((SELECT id FROM counties WHERE code = 'KSM'), 'Kisumu West'),
((SELECT id FROM counties WHERE code = 'KSM'), 'Seme'),
((SELECT id FROM counties WHERE code = 'KSM'), 'Nyando'),
((SELECT id FROM counties WHERE code = 'KSM'), 'Muhoroni'),
((SELECT id FROM counties WHERE code = 'KSM'), 'Nyakach');

-- Machakos County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'MCK'), 'Machakos Town'),
((SELECT id FROM counties WHERE code = 'MCK'), 'Athi River'),
((SELECT id FROM counties WHERE code = 'MCK'), 'Mavoko'),
((SELECT id FROM counties WHERE code = 'MCK'), 'Kathiani'),
((SELECT id FROM counties WHERE code = 'MCK'), 'Matungulu'),
((SELECT id FROM counties WHERE code = 'MCK'), 'Yatta'),
((SELECT id FROM counties WHERE code = 'MCK'), 'Kangundo'),
((SELECT id FROM counties WHERE code = 'MCK'), 'Masinga');

-- Kajiado County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'KJD'), 'Kajiado North'),
((SELECT id FROM counties WHERE code = 'KJD'), 'Kajiado Central'),
((SELECT id FROM counties WHERE code = 'KJD'), 'Kajiado East'),
((SELECT id FROM counties WHERE code = 'KJD'), 'Kajiado West'),
((SELECT id FROM counties WHERE code = 'KJD'), 'Kajiado South');

-- Meru County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'MRU'), 'Imenti North'),
((SELECT id FROM counties WHERE code = 'MRU'), 'Imenti South'),
((SELECT id FROM counties WHERE code = 'MRU'), 'Imenti Central'),
((SELECT id FROM counties WHERE code = 'MRU'), 'Buuri'),
((SELECT id FROM counties WHERE code = 'MRU'), 'Igembe South'),
((SELECT id FROM counties WHERE code = 'MRU'), 'Igembe Central'),
((SELECT id FROM counties WHERE code = 'MRU'), 'Igembe North'),
((SELECT id FROM counties WHERE code = 'MRU'), 'Tigania West'),
((SELECT id FROM counties WHERE code = 'MRU'), 'Tigania East');

-- Nyeri County Sub-Counties
INSERT INTO sub_counties (county_id, name) VALUES
((SELECT id FROM counties WHERE code = 'NYR'), 'Tetu'),
((SELECT id FROM counties WHERE code = 'NYR'), 'Kieni'),
((SELECT id FROM counties WHERE code = 'NYR'), 'Mathira'),
((SELECT id FROM counties WHERE code = 'NYR'), 'Othaya'),
((SELECT id FROM counties WHERE code = 'NYR'), 'Mukurweini'),
((SELECT id FROM counties WHERE code = 'NYR'), 'Nyeri Town');

