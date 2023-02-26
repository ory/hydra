  

update ory_keto.keto_relation_tuples set crdb_region = 'gcp-europe-west3' where crdb_region != 'gcp-europe-west3';
update ory_keto.keto_uuid_mappings set crdb_region = 'gcp-europe-west3' where crdb_region != 'gcp-europe-west3';
update ory_backoffice.active_projects set crdb_region = 'gcp-europe-west3' where crdb_region != 'gcp-europe-west3';


SELECT count(id), crdb_region from ory_backoffice.active_projects group by crdb_region;



SELECT id FROM ory_keto.keto_uuid_mappings WHERE string_representation='2339d03a-41d5-4c76-9dc2-dc69a5c56f2c'; -- object
SELECT id FROM ory_keto.keto_uuid_mappings WHERE string_representation='a9f362fa-72dc-4067-8c49-46e0f533fb42'; -- subject


SELECT keto_relation_tuples.commit_time, keto_relation_tuples.namespace, keto_relation_tuples.nid, keto_relation_tuples.object, keto_relation_tuples.relation, keto_relation_tuples.shard_id, keto_relation_tuples.subject_id, keto_relation_tuples.subject_set_namespace, keto_relation_tuples.subject_set_object, keto_relation_tuples.subject_set_relation FROM ory_keto.keto_relation_tuples AS keto_relation_tuples WHERE nid = '26f8b768-4797-4cc4-987f-b2aad16ca9f4' AND shard_id > '00000000-0000-0000-0000-000000000000' AND namespace = 'projects' AND object = '05455bbf-9413-5b37-b5fb-04cb4a259dc4' AND relation = 'owner' AND subject_id = '7479d00e-79b3-52b4-b335-538d4876dc57' AND subject_set_namespace IS NULL AND subject_set_object IS NULL AND subject_set_relation IS NULL ORDER BY shard_id, nid LIMIT 2;

SELECT keto_relation_tuples.commit_time, keto_relation_tuples.namespace, keto_relation_tuples.nid, keto_relation_tuples.object, keto_relation_tuples.relation, keto_relation_tuples.shard_id, keto_relation_tuples.subject_id, keto_relation_tuples.subject_set_namespace, keto_relation_tuples.subject_set_object, keto_relation_tuples.subject_set_relation FROM ory_keto_single_region.keto_relation_tuples AS keto_relation_tuples WHERE nid = '26f8b768-4797-4cc4-987f-b2aad16ca9f4' AND shard_id > '00000000-0000-0000-0000-000000000000' AND namespace = 'projects' AND object = '05455bbf-9413-5b37-b5fb-04cb4a259dc4' AND relation = 'owner' AND subject_id = '7479d00e-79b3-52b4-b335-538d4876dc57' AND subject_set_namespace IS NULL AND subject_set_object IS NULL AND subject_set_relation IS NULL ORDER BY shard_id, nid LIMIT 2;



ALTER DATABASE movr SET SECONDARY REGION



SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='353749b2-98b8-4f23-8c97-b457c6d4fca9';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='a9f362fa-72dc-4067-8c49-46e0f533fb42';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='00075931-1514-45f6-bd14-7916a8e307a0';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='000a7d06-d365-4dae-b74c-f3560859b4ce' LIMIT 10;
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='000c018d-8c35-409b-8336-7b3b5796b5bf';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='0010c7c8-ee21-483f-818f-aed69614d3d2';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='001334af-662c-4e30-8af1-3ec430daed01';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='0016d4dc-55a3-43b1-9d9a-767e15afce07';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='001b924f-986a-4766-b200-ac03079c58cd';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='001cfa5d-33ee-46d7-8eaf-9b9cf17e3ac9';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='001fb4d2-2490-43a7-a58d-c4eb596ae7b7';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='0020af0b-53f4-4a42-8a50-5629097ee83b';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='0022529c-3a2f-444f-a4af-4c6f7416556f';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='00251d75-14b2-4ef5-8d1a-b385c663c406';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='002cc8b6-a7c8-4a36-b372-daba738cf5da';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='002e5796-f4eb-4a94-94c0-69d7e05931e1';
SELECT * FROM ory_backoffice.active_projects WHERE identity_id=gen_random_uuid();
SELECT * FROM ory_backoffice.active_projects WHERE identity_id='002fe011-de63-4973-a434-1c751be44dd4';
