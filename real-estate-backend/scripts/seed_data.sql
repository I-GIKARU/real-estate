-- Kenyan Real Estate Database Seed Data
-- This script populates the database with Kenya counties and major sub-counties

-- Counties (47 counties in Kenya)
INSERT INTO counties (name, code, created_at) VALUES
('Mombasa', '001', NOW()),
('Kwale', '002', NOW()),
('Kilifi', '003', NOW()),
('Tana River', '004', NOW()),
('Lamu', '005', NOW()),
('Taita-Taveta', '006', NOW()),
('Garissa', '007', NOW()),
('Wajir', '008', NOW()),
('Mandera', '009', NOW()),
('Marsabit', '010', NOW()),
('Isiolo', '011', NOW()),
('Meru', '012', NOW()),
('Tharaka-Nithi', '013', NOW()),
('Embu', '014', NOW()),
('Kitui', '015', NOW()),
('Machakos', '016', NOW()),
('Makueni', '017', NOW()),
('Nyandarua', '018', NOW()),
('Nyeri', '019', NOW()),
('Kirinyaga', '020', NOW()),
('Murang''a', '021', NOW()),
('Kiambu', '022', NOW()),
('Turkana', '023', NOW()),
('West Pokot', '024', NOW()),
('Samburu', '025', NOW()),
('Trans-Nzoia', '026', NOW()),
('Uasin Gishu', '027', NOW()),
('Elgeyo-Marakwet', '028', NOW()),
('Nandi', '029', NOW()),
('Baringo', '030', NOW()),
('Laikipia', '031', NOW()),
('Nakuru', '032', NOW()),
('Narok', '033', NOW()),
('Kajiado', '034', NOW()),
('Kericho', '035', NOW()),
('Bomet', '036', NOW()),
('Kakamega', '037', NOW()),
('Vihiga', '038', NOW()),
('Bungoma', '039', NOW()),
('Busia', '040', NOW()),
('Siaya', '041', NOW()),
('Kisumu', '042', NOW()),
('Homa Bay', '043', NOW()),
('Migori', '044', NOW()),
('Kisii', '045', NOW()),
('Nyamira', '046', NOW()),
('Nairobi', '047', NOW())
ON CONFLICT (code) DO NOTHING;

-- Sub-counties for Nairobi (most relevant for real estate)
INSERT INTO sub_counties (county_id, name, created_at) VALUES
((SELECT id FROM counties WHERE code = '047'), 'Westlands', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Dagoretti North', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Dagoretti South', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Langata', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Kibra', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Roysambu', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Kasarani', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Ruaraka', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Embakasi South', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Embakasi North', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Embakasi Central', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Embakasi East', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Embakasi West', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Makadara', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Kamukunji', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Starehe', NOW()),
((SELECT id FROM counties WHERE code = '047'), 'Mathare', NOW());

-- Sub-counties for Kiambu (popular residential area)
INSERT INTO sub_counties (county_id, name, created_at) VALUES
((SELECT id FROM counties WHERE code = '022'), 'Thika Town', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Ruiru', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Juja', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Gatundu South', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Gatundu North', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Githunguri', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Kiambu', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Kiambaa', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Kabete', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Kikuyu', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Limuru', NOW()),
((SELECT id FROM counties WHERE code = '022'), 'Lari', NOW());

-- Sub-counties for Mombasa
INSERT INTO sub_counties (county_id, name, created_at) VALUES
((SELECT id FROM counties WHERE code = '001'), 'Changamwe', NOW()),
((SELECT id FROM counties WHERE code = '001'), 'Jomba', NOW()),
((SELECT id FROM counties WHERE code = '001'), 'Kisauni', NOW()),
((SELECT id FROM counties WHERE code = '001'), 'Nyali', NOW()),
((SELECT id FROM counties WHERE code = '001'), 'Likoni', NOW()),
((SELECT id FROM counties WHERE code = '001'), 'Mvita', NOW());

-- Sub-counties for Nakuru
INSERT INTO sub_counties (county_id, name, created_at) VALUES
((SELECT id FROM counties WHERE code = '032'), 'Nakuru Town East', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Nakuru Town West', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Bahati', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Subukia', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Rongai', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Kuresoi South', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Kuresoi North', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Gilgil', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Naivasha', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Molo', NOW()),
((SELECT id FROM counties WHERE code = '032'), 'Njoro', NOW());

-- Sub-counties for Kajiado
INSERT INTO sub_counties (county_id, name, created_at) VALUES
((SELECT id FROM counties WHERE code = '034'), 'Kajiado North', NOW()),
((SELECT id FROM counties WHERE code = '034'), 'Kajiado Central', NOW()),
((SELECT id FROM counties WHERE code = '034'), 'Kajiado East', NOW()),
((SELECT id FROM counties WHERE code = '034'), 'Kajiado West', NOW()),
((SELECT id FROM counties WHERE code = '034'), 'Magadi', NOW());
