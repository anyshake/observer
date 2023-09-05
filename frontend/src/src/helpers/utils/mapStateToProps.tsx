import { ReduxStore } from "../../config/store";

const mapStateToProps = (state: ReduxStore) => {
    return { ...state };
};

export default mapStateToProps;
