import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../providers/label_notifier_provider.dart';
import '../../themes/colors.dart';
import '../../../../domain/kanban/models/card_item.dart';
import '../../../kanban/widgets/card_details_widgets/cart_details_widget.dart';
import 'card_item_content.dart';

class CardItemDesktop extends ConsumerStatefulWidget {
  final CardItem card;
  const CardItemDesktop({super.key, required this.card});

  @override
  ConsumerState<CardItemDesktop> createState() => _CardItemDesktopState();
}

class _CardItemDesktopState extends ConsumerState<CardItemDesktop> {
  bool isHovered = false;
  
  @override
  Widget build(BuildContext context) {
    final labels = ref.watch(labelNotifierProvider).labels
        .where((label) => label.cardId == widget.card.id)
        .toList();

    Widget content = CardItemContent(card: widget.card, labels: labels);

    return GestureDetector(
      onTap: () => _openDialog(context),
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        onEnter: (_) => setState(() => isHovered = true),
        onExit: (_) => setState(() => isHovered = false),
        child: Container(
          width: double.infinity,
          padding: const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color: isHovered ? AppColors.white31 : AppColors.white15,
            borderRadius: BorderRadius.circular(15),
          ),
          child: content,
        ),
      ),
    );
  }

  void _openDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => Dialog(
        insetPadding: const EdgeInsets.all(16),
        backgroundColor: Colors.transparent,
        child: CartDetailsWidget(card: widget.card),
      ),
    );
  }
}