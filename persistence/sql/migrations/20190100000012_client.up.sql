ALTER TABLE hydra_client ADD created_at timestamp NOT NULL DEFAULT now();
ALTER TABLE hydra_client ADD updated_at timestamp NOT NULL DEFAULT now();
