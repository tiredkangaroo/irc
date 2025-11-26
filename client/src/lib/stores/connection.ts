import { writable, type Writable } from "svelte/store";
import type { Connection } from "../../types";

// connection will be of type Connection if there is connection info saved.
// undefined means the connection info is being fetched.
// null means there is no connection info saved.
export const connection: Writable<Connection | undefined | null> =
  writable(null);
