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

    // vector <Pos> bazookas;
    // vector <Pos> guns;
    // vector <Pos> moneys;
    // vector <Pos> foods;
    // bool start = true;

    const vector<Dir> dirs = {Up, Down, Left, Right};
    // adalt, abaix, esquerra, dreta
    const vector<vector<int>> mov = {{-1, 0}, {1, 0}, {0, -1}, {0, 1}};

    map<Pos, Pos> Path;

    Pos BFS(Pos p, char c)
    {
        vector<vector<bool>> visitats(board_cols(), vector<bool>(board_rows(), false));
        queue<Pos> cola;
        cola.push(p);
        while (!cola.empty())
        {
            // p has to be the node visited
            Pos pact = cola.front();
            cola.pop();
    
            Cell c = cell(pact);
            visitats[pact.i][pact.j] = true;

            if (c.weapon == Bazooka)
                return pact;
            else if (c.weapon == Gun)
                return pact;

            for (int i = 0; i < 4; ++i)
            {
                Pos ptemp = pact;
                ptemp.i += mov[i][0];
                ptemp.j += mov[i][1];
                Cell c = cell(ptemp);
                if (pos_ok(ptemp) and c.type == Street and not visitats[ptemp.i][ptemp.j])
                {
                    Path[ptemp] = pact;
                    cola.push(ptemp);
                }
            }
        }
        return Pos(-1,-1);
    }
    Dir Pos2dir (Pos p) {
        if (p == Pos(-1,0)) return Up;
        else if (p == Pos(1,0)) return Down;
        else if (p == Pos(0,-1)) return Left;
        else if (p == Pos(0,1)) return Right;
        else return Up;
    }
    // Invariant: l'objecte ha d'existir en el mapa en el moment de la crida i l'hem trobat amb el BFS
    Dir direccionBFS(Pos pciti, Pos pobj) {
        Pos ptemp = pobj;
        Pos ptemp2 = Path.at(ptemp);
        while (pciti != ptemp) {
            ptemp = Path.at(ptemp);
            ptemp2 = Path.at(ptemp2);
        }
        Pos r; 
        r.i = ptemp.i - pciti.i;
        r.j = ptemp.j - pciti.j;
        cerr << r;
        Dir dr = Pos2dir(r);
        return dr;
    }

    /**
   * Play method, invoked once per each round.
   */
    virtual void play()
    {
        // if (start) {
        //     start = false;
        //     bazookas.resize(num_ini_bazookas());
        //     guns.resize(num_ini_guns());
        //     moneys.resize(num_ini_money());
        //     foods.resize(num_ini_food());
        //     for (int i = 0; i < num_ini_bazookas(); ++i) {
        //         bazookas[i].i = -1;
        //         bazookas[i].j = -1;
        //     }
        //     for (int i = 0; i < num_ini_guns(); ++i) {
        //         guns[i].i = -1;
        //         guns[i].j = -1;
        //     }
        //     for (int i = 0; i < num_ini_money(); ++i) {
        //         moneys[i].i = -1;
        //         moneys[i].j = -1;
        //     }
        //     for (int i = 0; i < num_ini_food(); ++i) {
        //         foods[i].i = -1;
        //         foods[i].j = -1;
        //     }
        // }

        // n = board max column, m = board max row
        // int n = board_cols();
        // int m = board_rows();
        // for (int i = 0; i < n; ++i)
        // {
        //     for (int j = 0; j < m; ++j)
        //     {
        //         Cell casilla = cell(i, j);
        //         if (casilla.weapon == Bazooka) {
        //
        //         }
        //         cerr << casilla.type << ' '; // cout << street or building
        //     }
        //     cerr << endl;
        // }

        // INICIALITZAR INVARIANTS
        Path.clear();
        if (is_day())
        {
            vector<int> w = warriors(me());
            for (int id : w)
            { // iterate over warriors
                Pos pciti = citizen(id).pos;
                cerr << "citizen " << id << " in position: " << pciti << endl;
                char c = 'b'; // busquem bazokas
                Pos pobj = BFS(pciti,c);
                Dir d = direccionBFS(pciti,pobj);
                //Dir d = Down; // prova
                if (pos_ok(pciti + d))
                {
                    move(id, d);
                }
            }

            // At day take care of builders
            vector<int> b = builders(me());
            for (int id : b)
            { // iterate over builders
                Pos p = citizen(id).pos;
                cerr << "citizen " << id << " in position: " << p << endl;
                // BFS(p);

                if (pos_ok(p + Down))
                {
                    move(id, Down);
                }
            }
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
