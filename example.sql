-- Anti Pattern 7: Convert Dynamic Predicates into Static
SELECT
    col1, col2 
FROM 
    table1 
WHERE 
    col1 > 0
    AND col3 in (select col3 from table2)
    AND col2 = 1;

-- SELECT
--     col1, col2 
-- FROM 
--     table1 
-- WHERE 
--     col3 in (select col3 from table2);


-- SELECT
--     col1, col2 
-- FROM 
--     table1 
-- WHERE 
--     col3 in (select col3 from table2)
--     AND col1>0;
