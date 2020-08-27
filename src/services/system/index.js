import reducer from "./reducer";
import saga from "./sagas";

// Service > system

export const mountPoint = "system";

export default {
    mountPoint,
    reducer,
    saga
};
