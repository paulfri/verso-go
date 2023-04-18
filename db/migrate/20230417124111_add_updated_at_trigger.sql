-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION touch_updated_at()   
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;   
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION touch_updated_at();
-- +goose StatementEnd
