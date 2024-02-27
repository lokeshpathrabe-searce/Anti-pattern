-- Anti Pattern 2: SEMI-JOIN without aggregation
SELECT 
    t1.col1 
FROM 
    `project.dataset.table1` t1 
WHERE 
    t1.col1 > 0 
    AND t1.col2 IN (SELECT col2 FROM `project.dataset.table2`) 
    AND t1.col3 IN (SELECT col3 FROM `project.dataset.table3` ) 
    AND t1.col4 IN (1, 2, 3, 4);
-- SELECT 
--    t1.col1 
-- FROM 
--    `project.dataset.table1` t1 
-- WHERE 
--     t1.col2 not in (select col2,col3 from `project.dataset.table2`);
-- SELECT 
--    t1.col1 
-- FROM 
--    `project.dataset.table1` t1 
-- WHERE 
--     t1.col2 not in (select col2, id from `project.dataset.table2`);
-- SELECT
--     col1, col2 
-- FROM 
--     table1 
-- WHERE 
--     col1 > 0
--     AND col3 IN (SELECT col3 FROM table2)
--     AND col2 = 1;