import React, { FC, useState } from 'react';
import { css } from '@emotion/css';
import { GrafanaTheme } from '@grafana/data';
import { Button, useStyles } from '@grafana/ui';
import { AmRouteReceiver, FormAmRoute } from '../../types/amroutes';
import { emptyRoute } from '../../utils/amroutes';
import { EmptyArea } from '../EmptyArea';
import { AmRoutesTable } from './AmRoutesTable';

export interface AmSpecificRoutingProps {
  onChange: (routes: FormAmRoute) => void;
  receivers: AmRouteReceiver[];
  routes: FormAmRoute;
}

export const AmSpecificRouting: FC<AmSpecificRoutingProps> = ({ onChange, receivers, routes }) => {
  const [actualRoutes, setActualRoutes] = useState(routes.routes);
  const [isAddMode, setIsAddMode] = useState(false);

  const styles = useStyles(getStyles);

  const addNewRoute = () => {
    setIsAddMode(true);
    setActualRoutes((actualRoutes) => [...actualRoutes, emptyRoute]);
  };

  return (
    <div className={styles.container}>
      <h5>Specific routing</h5>
      <p>Send specific alerts to chosen channels, based on matching criteria</p>
      {actualRoutes.length > 0 ? (
        <>
          {!isAddMode && (
            <Button className={styles.addMatcherBtn} icon="plus" onClick={addNewRoute} type="button">
              New policy
            </Button>
          )}
          <AmRoutesTable
            isAddMode={isAddMode}
            onChange={(newRoutes) => {
              onChange({
                ...routes,
                routes: newRoutes,
              });

              if (isAddMode) {
                setIsAddMode(false);
              }
            }}
            receivers={receivers}
            routes={actualRoutes}
          />
        </>
      ) : (
        <EmptyArea
          buttonIcon="plus"
          buttonLabel="New specific policy"
          onButtonClick={addNewRoute}
          text="You haven't created any specific policies yet."
        />
      )}
    </div>
  );
};

const getStyles = (_theme: GrafanaTheme) => {
  return {
    container: css`
      display: flex;
      flex-flow: column nowrap;
    `,
    addMatcherBtn: css`
      align-self: flex-end;
      margin-bottom: 28px;
    `,
  };
};
