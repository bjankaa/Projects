package MenuPack;

import ScorePack.ScoreFrame;
import GamePack.GameFrame;
import javax.swing.*;
import javax.swing.border.EmptyBorder;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.IOException;

public class MenuFrame extends JFrame{

    private void initComponents(){

        //Gomboknak
        JPanel menu = new JPanel();
        menu.setLayout(new GridLayout(3,1));
        menu.setBorder (new EmptyBorder(new Insets(70,200,250,200)));

        this.setLayout(new BorderLayout());
        this.add(menu, BorderLayout.CENTER);

        //Felirat
        JLabel renju = new JLabel("Renju"); ////// EDIT
        renju.setHorizontalAlignment(JLabel.CENTER);
        renju.setPreferredSize(new Dimension(100, 150));
        renju.setFont(new Font("Serif", Font.BOLD, 50));
        this.add(renju, BorderLayout.NORTH);

        ///////Menu panel resze//////////

        // Jatek menupont init
        JButton jatek = new JButton("Játék");
        jatek.setFont(new Font("Serif", Font.PLAIN, 30));
        menu.add(jatek);
        jatek.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                setVisible(false);
                GameFrame gf = new GameFrame();
                gf.run();
            }
        });

        // Eredmeny menupont init
        JButton eredmeny = new JButton("Eredmények");
        eredmeny.setFont(new Font("Serif", Font.PLAIN, 30));
        menu.add(eredmeny);
        eredmeny.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                setVisible(false);
                ScoreFrame sf = null;
                try {
                    sf = new ScoreFrame();
                } catch (IOException ex) {
                    ex.printStackTrace();
                }
                sf.run();
            }
        });

        // Szabaly menupont init
        JButton szabalyok = new JButton ("Szabályok");
        szabalyok.setFont(new Font("Serif", Font.PLAIN, 30));
        menu.add(szabalyok);
        szabalyok.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                setVisible(false);
                RuleFrame rf = new RuleFrame();
                rf.run();
            }
        });
    }

    // MenuFrame letrehozzasa
    public MenuFrame() {
        super("RENJU");
        setDefaultCloseOperation(EXIT_ON_CLOSE);
        setSize(800, 700);
        setLocationRelativeTo(null);
        setResizable(false);
        initComponents();
    }

    public void run(){
        this.setVisible(true);
    }

    public static void main(String[] args) {
        MenuFrame mf = new MenuFrame();
        mf.run();
    }
}