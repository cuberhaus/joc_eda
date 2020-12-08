#include "Player.hh"

/**
 * Write the name of your player and save this file
 * with the same name and .cc extension.
 */
#define PLAYER_NAME Null

void BFS(Pos p, int dest, int v, int)
{
}

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

    /**
   * Play method, invoked once per each round.
   */
    virtual void play()
    {
        if (is_day())
        {
			vector<int> w = warriors(me());
			for (int id : w) { // iterate over warriors
                Pos p = citizen(id).pos;
                // BFS(p);
                Dir d = Down; // prova
                move(id, d);
            }

			// At day take care of builders
			vector<int> b = builders(me());
			for (int id : b) { // iterate over builders
                Pos p = citizen(id).pos;
                // BFS(p);
                Dir d = Down; // prova
                move(id, d);
            }

            // n = board max column, m = board max row
            int n = board_cols();
            int m = board_rows();
            for (int i = 0; i < n; ++i)
            {
                for (int j = 0; j < m; ++j)
                {
                    Cell casilla = cell(i, j);
                    cerr << casilla.type << ' '; // cout << street or building
                }
                cerr << endl;
            }
            cerr << "això es una prova";
        }
        else { // night time

        }
    }
};

/**
 * Do not modify the following line.
 */
RegisterPlayer(PLAYER_NAME);
