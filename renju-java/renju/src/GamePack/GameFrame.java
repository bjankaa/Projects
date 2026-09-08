package GamePack;

import MenuPack.MenuFrame;

import javax.swing.*;
import javax.swing.border.EmptyBorder;

import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.IOException;

import static javax.swing.WindowConstants.EXIT_ON_CLOSE;

public class GameFrame extends JFrame {

    private JButton[][] grid;
    public String[][] gombtomb;
    public WinFrameWhite wfw = new WinFrameWhite();
    public WinFrameBlack wfb = new WinFrameBlack();

    //Atlatszova teszi a gombjainkat (for aesthetic purposes)
    public void atlatszobe(JButton gomb) {
        gomb.setOpaque(false);
        gomb.setContentAreaFilled(false);
        gomb.setBorderPainted(false);
    }

    // Fekete kovet jeleniti meg amikor kattintunk
    public void feketekor(JButton gomb) {
        var fekete = new ImageIcon("fekete.png");
        var img = fekete.getImage();
        img = img.getScaledInstance(gomb.getHeight() - 8, gomb.getHeight() - 8, Image.SCALE_SMOOTH);
        fekete = new ImageIcon(img);
        gomb.setIcon(fekete);
    }

    // Feher kovet jeleniti meg amikor kattintunk
    public void feherkor(JButton gomb) {
        var feher = new ImageIcon("feher.png");
        var img = feher.getImage();
        img = img.getScaledInstance(gomb.getHeight() + 11, gomb.getHeight() + 11, Image.SCALE_SMOOTH);
        feher = new ImageIcon(img);
        gomb.setIcon(feher);
    }
    //Jatek ter inicializalasa
    private void initComponents(){
        this.setLayout(new BorderLayout());

        JLabel tabla = new JLabel();
        this.add(tabla, BorderLayout.CENTER);
        ImageIcon hatter = new ImageIcon("tabla.png");
        var kep = hatter.getImage();
        kep = kep.getScaledInstance(885, 727, Image.SCALE_SMOOTH);
        hatter = new ImageIcon(kep);
        tabla.setIcon(hatter);
        tabla.setLayout(new GridLayout(15, 15));

        //Gombok inicializalasa
        this.grid = new JButton[15][15];
        for (int x = 0; x < 15; x++) {
            for (int y = 0; y < 15; y++) {
                grid[x][y] = new JButton();

                tabla.add(grid[x][y]);
                atlatszobe(grid[x][y]);
            }
        }

        JPanel alsosav = new JPanel();
        this.add(alsosav, BorderLayout.SOUTH);

        //Smallframe
        JButton pause = new JButton("Pause");
        pause.setFont(new Font("Serif", Font.PLAIN, 20));
        alsosav.add(pause);
        SmallMenuFrame smf = new SmallMenuFrame();
        smf.setGameFrame(this);
        pause.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                setEnabled(false);
                smf.run();
            }
        });
    }

    // A ket dimenzios tomb inicializalasa, amiben taruljuk a tablankat.
    // Adunk neki egy keretet, ami ellenorozeseket szempontjabol hasznos. Ezert nagyobb is mint a jatek tablank
    public void initGombTomb() {
        this.gombtomb = new String[17][17];
        for (int i = 0; i < 17; i++) {
            for (int j = 0; j < 17; j++) {
                boolean top = (i == 0);
                boolean bottom = (i == 16);
                boolean left = (j == 0);
                boolean right = (j == 16);
                if (top || bottom || left || right) {
                    gombtomb[i][j] = "#";
                } else {
                    gombtomb[i][j] = "-";
                }
            }
        }
    }

    //Ez teszibe az elemeket abba a tombe ahol nyilvan tartjuk a tabla allasat( color : X feket, O feher)
    public void gombtombeElem(JButton gomb, String color) throws IOException {
        for (int x = 0; x < 15; x++) {
            for (int y = 0; y < 15; y++) {
                if (grid[x][y] == gomb) {
                    gombtomb[x + 1][y + 1] = color;
                    Game jtk = new Game();
                    jtk.Game(gombtomb, x + 1, y + 1, color, wfw, wfb);
                }
            }
        }
    }

    public GameFrame() {
        super("RENJU - Játék");
        setDefaultCloseOperation(EXIT_ON_CLOSE);
        setSize(900, 800);
        setLocationRelativeTo(null);
        setResizable(false);
        wfw.setGameFrame(this);
        wfb.setGameFrame(this);
        initComponents();
        initGombTomb();

        //Kozepso elem fekete, kotelezo fekete kezdes
        try {
            gombtombeElem(grid[7][7], "X");
        } catch (IOException e) {
            e.printStackTrace();
        }
        feketekor(grid[7][7]);

        //Gombnyomasok

        for (int x = 0; x < 15; x++) {
            for (int y = 0; y < 15; y++) {
                if (grid[x][y] == grid[7][7]) {
                    continue;
                }
                this.grid[x][y].addActionListener(new ActionListener() {
                    private static int lepes = 0;
                    @Override
                    public void actionPerformed(ActionEvent e) {
                        if (lepes % 2 == 0) {
                            feherkor((JButton) e.getSource());
                            try {
                                gombtombeElem((JButton) e.getSource(), "O");
                            } catch (IOException ex) {
                                ex.printStackTrace();
                            }
                            ((JButton) e.getSource()).removeActionListener(this);
                            lepes++;
                        } else if (lepes % 2 == 1) {
                            feketekor((JButton) e.getSource());
                            try {
                                gombtombeElem((JButton) e.getSource(), "X");
                            } catch (IOException ex) {
                                ex.printStackTrace();
                            }
                            ((JButton) e.getSource()).removeActionListener(this);
                            lepes++;
                        }
                    }
                });
            }
        }

    }
    public void run(){
        this.setVisible(true);
    }

    public void close(){
        this.dispose();
    }

}
