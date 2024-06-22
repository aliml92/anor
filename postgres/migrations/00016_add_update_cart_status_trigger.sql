-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_cart_status_on_order_creation()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE carts
    SET status = 'Completed'
    WHERE id = (SELECT cart_id FROM orders WHERE id = NEW.id);

    RETURN NEW;
END;
$$;

CREATE TRIGGER update_cart_status_trig
    AFTER INSERT ON orders
    FOR EACH ROW
EXECUTE FUNCTION update_cart_status_on_order_creation();
-- +goose StatementEnd
