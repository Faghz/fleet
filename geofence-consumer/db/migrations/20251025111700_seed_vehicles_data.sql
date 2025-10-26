-- migrate:up

-- Insert sample vehicles with UUID and Indonesian vehicle number format
INSERT INTO public.vehicle (entity_id, vehicle_id, vehicle_type, brand, model, year, status, created_by) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'B1234ABC', 'intercity_bus', 'Mercedes-Benz', 'OH 1626', 2022, 'active', 'system'),
('550e8400-e29b-41d4-a716-446655440002', 'B5678DEF', 'city_bus', 'Isuzu', 'Elf NLR', 2021, 'active', 'system'),
('550e8400-e29b-41d4-a716-446655440003', 'D9012GHI', 'minibus', 'Toyota', 'Hiace', 2023, 'active', 'system'),
('550e8400-e29b-41d4-a716-446655440004', 'B3456JKL', 'double_decker', 'Scania', 'K410', 2020, 'active', 'system'),
('550e8400-e29b-41d4-a716-446655440005', 'E7890MNO', 'intercity_bus', 'Hino', 'RK8', 2022, 'active', 'system');

-- migrate:down
DELETE FROM public.vehicle WHERE created_by = 'system' AND entity_id IN (
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440003',
    '550e8400-e29b-41d4-a716-446655440004',
    '550e8400-e29b-41d4-a716-446655440005'
);
