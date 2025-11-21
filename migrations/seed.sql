INSERT INTO public.places (name, address, point, rating_sum, rating_count, image_url) VALUES
  (
    'Pat''s King of Steaks',
    '1237 E Passyunk Ave, Philadelphia, PA 19147',
    ST_GeomFromText('POINT(-75.1581 39.9308)', 4326),
    47,   -- avg ~4.27
    11,
    ARRAY['https://example.com/pats1.jpg']
  ),
  (
    'Geno''s Steaks',
    '1219 S 9th St, Philadelphia, PA 19147',
    ST_GeomFromText('POINT(-75.1585 39.9299)', 4326),
    41,   -- avg ~3.73
    11,
    ARRAY['https://example.com/genos1.jpg']
  ),
  (
    'Jim''s Steaks South Street',
    '400 S St, Philadelphia, PA 19147',
    ST_GeomFromText('POINT(-75.1605 39.9483)', 4326),
    52,   -- avg ~4.33
    12,
    ARRAY['https://example.com/jims1.jpg']
  ),
  (
    'Tony Luke''s',
    '39 E Oregon Ave, Philadelphia, PA 19148',
    ST_GeomFromText('POINT(-75.1608 39.9293)', 4326),
    44,   -- avg ~4.0
    11,
    ARRAY['https://example.com/tonylukes1.jpg']
  ),
  (
    'Dalessandro''s Steaks & Hoagies',
    '600 Wendover St, Philadelphia, PA 19128',
    ST_GeomFromText('POINT(-75.2210 40.0184)', 4326),
    49,   -- avg ~4.45
    11,
    ARRAY['https://example.com/dalessandros1.jpg']
  );
