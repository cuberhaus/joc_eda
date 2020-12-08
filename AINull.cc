#include "Player.hh"

/**
 * Write the name of your player and save this file
 * with the same name and .cc extension.
 */
#define PLAYER_NAME Null

struct PLAYER_NAME : public Player
{

    /**
   * Factory: returns a new instance of this class.
   * Do not modify this function.
   */
    static Player *factory()
    {
        return new PLAYER_NAME;
    }

    /**
   * Types and attributes for your player can be defined here.
   */
    const vector<Dir> dirs = {Up, Down, Left, Right};
    // adalt, abaix, esquerra, dreta
    const vector<vector<int>> mov = {{-1, 0}, {1, 0}, {0, -1}, {0, 1}};

    void BFS(Pos p)
    {
        queue<Pos> cola;
        for (int i = 0; i < 4; ++i)
        {
            p.i += mov[i][0];
            p.j += mov[i][1];
            if (pos_ok(p))
            {
                cola.push(p);
            }
        }
    }

    /**
   * Play method, invoked once per each round.
   */
    virtual void play()
    {
        if (is_day())
        {
            vector<int> w = warriors(me());
            for (int id : w)
            { // iterate over warriors
                Pos p = citizen(id).pos;
                cerr << "citizen " << id << " in position: " << p << endl;
                // BFS(p);
                //Dir d = Down; // prova
                if (pos_ok(p+Down)) {
                    move(id, Down);
                }
            }

            // At day take care of builders
            vector<int> b = builders(me());
            for (int id : b)
            { // iterate over builders
                Pos p = citizen(id).pos;
                cerr << "citizen " << id << " in position: " << p << endl;
                // BFS(p);

                if (pos_ok(p+Down)) {
                    move(id, Down);
                }
            }

            // n = board max column, m = board max row
            //int n = board_cols();
            //int m = board_rows();
            //for (int i = 0; i < n; ++i)
            //{
            //    for (int j = 0; j < m; ++j)
            //    {
            //        Cell casilla = cell(i, j);
            //        cerr << casilla.type << ' '; // cout << street or building
            //    }
            //    cerr << endl;
            //}
            //cerr << "això es una prova";
        }
        else
        { // night time
        }
    }
};

/**
 * Do not modify the following line.
 */
RegisterPlayer(PLAYER_NAME);
