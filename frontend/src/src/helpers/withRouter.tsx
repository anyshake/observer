import { useNavigate, useLocation } from "react-router-dom";
import { ComponentType, JSXElementConstructor } from "react";

const withRouter = <P extends {}>(
    Component: ComponentType<P>
): JSXElementConstructor<
    P & {
        history: ReturnType<typeof useNavigate>;
        location: ReturnType<typeof useLocation>;
    }
> => {
    return (props: P) => {
        const history = useNavigate();
        const location = useLocation();
        return <Component {...props} history={history} location={location} />;
    };
};

export default withRouter;
