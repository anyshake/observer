import { useNavigate, useLocation } from "react-router-dom";

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name withRouter - 为 Component 注入属性
 * @description 为 Component 注入 history 和 location 属性
 * @param {JSX.Element} Component - React Component
 * @returns {JSX.Element} 返回新 React Component
 */
const withRouter = (Component) => (props) => {
    const history = useNavigate();
    const location = useLocation();
    return <Component {...props} history={history} location={location} />;
};

export default withRouter;
