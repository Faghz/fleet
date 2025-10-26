-- migrate:up

-- Insert sample points of interest
INSERT INTO public.point_of_interest (id, name, latitude, longitude, description, created_by, created_at) VALUES
(1, 'National Monument (Monas)', -6.1753924, 106.8271528, 'Iconic monument in Jakarta, Indonesia', 'system', NOW()),
(2, 'Istiqlal Mosque', -6.169856, 106.830759, 'Largest mosque in Southeast Asia', 'system', NOW()),
(3, 'Borobudur Temple', -7.6079, 110.2038, 'Ancient Buddhist temple in Magelang, Indonesia', 'system', NOW()),
(4, 'Jakarta Cathedral', -6.1690, 106.8330, 'Gothic-style cathedral in Jakarta', 'system', NOW()),
(5, 'Taman Mini Indonesia Indah', -6.3024, 106.8952, 'Cultural park showcasing Indonesian diversity', 'system', NOW()),
(6, 'Grand Indonesia Mall', -6.1951, 106.8227, 'Luxury shopping mall in Jakarta', 'system', NOW()),
(7, 'Ancol Dreamland', -6.1173, 106.8584, 'Recreational area with beach and amusement park in Jakarta', 'system', NOW()),
(8, 'Mount Bromo', -7.9425, 112.9533, 'Active volcano and popular tourist destination in East Java', 'system', NOW()),
(9, 'Kuta Beach', -8.7203, 115.1671, 'Famous beach in Bali known for surfing', 'system', NOW()),
(10, 'Prambanan Temple', -7.7520, 110.4915, 'Hindu temple complex in Yogyakarta, UNESCO World Heritage Site', 'system', NOW());

-- migrate:down
DELETE FROM public.point_of_interest WHERE id in (1,2,3,4,5,6,7,8,9,10);