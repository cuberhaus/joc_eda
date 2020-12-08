#include "Player.hh"

/**
 * Write the name of your player and save this file
 * with the same name and .cc extension.
 */
#define PLAYER_NAME Null

// BFS(vector<vector<int>> adj, int src, int dest, int v, int)
// {
// 
// }

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
        //// n = board max column, m = board max row
        int x;
        for (int i = 0; i < 10; ++i) {
            for (int j = 0; j < 10; ++j) {
                x = x+1;
            }
        }
        cerr << x;
        //Cell = cell(0, 0);
        cerr << "això es una prova";
    }
};

/**
 * Do not modify the following line.
 */
RegisterPlayer(PLAYER_NAME);
